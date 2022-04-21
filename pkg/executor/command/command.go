package command

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"os/user"
	"regexp"
	"rexagent/pkg/log"
	"rexagent/pkg/task"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type CommandExecutor struct {
	id   string
	name string

	taskProcesses map[int]*exec.Cmd
}

func (e *CommandExecutor) Execute(t *task.CommandJob, ctx *task.Context) (outBuf []byte, err error) {
	// 对job进行校验
	err = e.Validate(t)
	if err != nil {
		return nil, err
	}
	//必要的时候初始化上下文Context
	if ctx == nil {
		ctx = &task.Context{}
	}
	//初始化上下文的参数
	if t.Parameters != nil {
		for _, p := range t.Parameters {
			ctx.SetVar(p.Name, p.Value)
		}
	}



	cmd, out, err := e.buildCommand(t, ctx)
	if err != nil {
		return nil, err
	}
	if out != nil {
		defer out.Close()
	}

	outBuf, err = e.executeCommand(t, cmd, out)

	return outBuf, err

}

func (e *CommandExecutor) Validate(t *task.CommandJob) error {
	if t == nil {
		return errors.New("the job is null")
	}
	if t.Code == "" {
		return errors.New("the code/script is null")
	}

	if t.Async == 1 && t.CallbackApi == "" {
		return errors.New("the async call need the callback api")
	}

	return nil
}

func (e *CommandExecutor) buildCommand(t *task.CommandJob, ctx *task.Context) (cmd *exec.Cmd, out io.ReadCloser, err error) {
	var name string
	var args []string

	code := strings.TrimSpace(t.Code)
	switch t.Lang {
	case "bash", "":
		name, args = getCmdBashArgs(code, ctx)
	case "exec":
		var exists bool
		name, args, exists = getCmdExecArgs(code, ctx)
		if !exists {
			return nil, nil, errors.New("param missing")
		}
	default:
		return nil, nil, errors.New("unknown language")
	}
	cmd = exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.Dir = "/"

	if t.User != "" {
		err = setCmdUser(cmd, t.User)
		if err != nil {
			return nil, nil, errors.New("set user failed")
		}
	}
	cmd.Stdin = nil
	cmd.Stderr = nil
	if t.Background == 1 {
		cmd.SysProcAttr.Setsid = true
		cmd.SysProcAttr.Foreground = false
		cmd.Stdout = nil
		cmd.Stdin = nil
	} else {
		cmd.SysProcAttr.Setpgid = true
		cmd.SysProcAttr.Pgid = 0
		out, err = cmd.StdoutPipe()
		if err != nil {
			return nil, nil, errors.New("pipe stdout failed")
		}
	}
	return cmd, out, nil

}

func (e *CommandExecutor) executeCommand(t *task.CommandJob, cmd *exec.Cmd, out io.ReadCloser) (outBuf []byte, err error) {
	err = cmd.Start()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("execution err: %s", err))
	}
	log.Info("process %d running", cmd.Process.Pid)
	if t.Background == 1 {
		go func() {
			err = cmd.Wait()
			if err != nil {
				log.Error("background process %d ended with erro %d", cmd.Process.Pid, err.Error())
			}
		}()
	} else {
		ch := make(chan error, 1)
		go func() {
			if out != nil {
				if outBuf, err = ioutil.ReadAll(out); err != nil {
					ch <- err
					_ = cmd.Wait()
					return
				}
				if err = cmd.Wait(); err != nil {
					ch <- errors.New(fmt.Sprintf("execution error %s, log: %s", err, string(outBuf)))
					return
				}
				ch <- nil
				return
			}
		}()
		timeout := time.Duration(t.Timeout)
		select {
		case err = <-ch:
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error : %s", err.Error()))
			}
		case <-time.After(timeout * time.Second):
			_ = cmd.Process.Kill()
			return nil, errors.New("timeout")
		}
	}

	return

}

func getCmdBashArgs(code string, context *task.Context) (string, []string) {
	return "bash", []string{"-c", code}
}

var varExpr, _ = regexp.Compile(`^[a-zA-Z_]\w*(?:\.[a-zA-Z_]\w*)?$`)
var argRe, _ = regexp.Compile(`("[^"]*"|'[^']*'|[^\s"']+)`)

func getCmdExecArgs(code string, ctx *task.Context) (string, []string, bool) {
	argsMatches := argRe.FindAllStringSubmatch(code, -1)
	args := make([]string, 0, 4)
	var exists bool
	for i := 0; i < len(argsMatches); i++ {
		arg := argsMatches[i][1]
		if arg[0] == '\'' || arg[0] == '"' {
			arg = arg[1 : len(arg)-1]
		}
		arg, exists = replaceCmdParams(arg, ctx)
		if !exists {
			return "", nil, false
		}
		args = append(args, arg)
	}

	return args[0], args[1:], true
}

func replaceCmdParams(arg string, context *task.Context) (string, bool) {
	return VarExpand(arg, func(s string) string {
		ret, _ := context.GetVar(s)
		return ret
	})
}

func VarExpand(s string, replace func(string) string) (string, bool) {
	const maxDepth = 10
	stack := make([][]byte, maxDepth)
	stack[0] = make([]byte, 0, len(s))
	for i := 1; i < len(stack); i++ {
		stack[i] = make([]byte, 0, 8)
	}
	sp := 0
	for i := 0; i < len(s); i++ {
		if i < len(s)-1 && s[i:i+2] == "${" {
			sp++
			if sp == maxDepth {
				return "", false
			}
			i++
		} else if s[i] == '}' {
			if sp < 1 {
				return "", false
			}
			n := stack[sp]
			if !varExpr.Match(n) {
				return "", false
			}

			if sp == 1 {
				stack[0] = append(stack[0], []byte(replace(string(n)))...)
			} else {
				stack[sp-1] = append(stack[sp-1], []byte(string(n))...)
			}
			stack[sp] = stack[sp][:0]
			sp--
		} else {
			stack[sp] = append(stack[sp], s[i])
		}
	}
	if sp != 0 {
		return "", false
	}
	return string(stack[sp]), true
}

func setCmdUser(cmd *exec.Cmd, username string) error {
	sysUser, err := user.Lookup(username)
	if err != nil {
		return err
	}
	uid, err := strconv.Atoi(sysUser.Uid)
	if err != nil {
		return err
	}
	gid, err := strconv.Atoi(sysUser.Gid)
	if err != nil {
		return err
	}
	cred := syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}

	cmd.SysProcAttr.Credential = &cred
	return nil
}
