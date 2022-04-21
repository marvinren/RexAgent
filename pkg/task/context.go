package task

import (
	"rexagent/pkg/variable"
	"sync"
)

type Context struct {
	//参数（包括入参，任务之间传递的参数）
	Args map[string]string
	//全局变量
	GlobalVars map[string]string
	//Args的锁
	argsLock sync.Mutex
}

func NewContext(initParams string) *Context {

	return &Context{
		Args: variable.GetArgumentsFromStr(initParams),
	}
}

func (c *Context) GetVar(name string) (string, bool) {
	//读取变量需要有锁
	c.argsLock.Lock()
	defer c.argsLock.Unlock()
	v, ok := c.Args[name]

	if !ok {
		//如果在上下文参数中无法找到变量，从全局变量中寻找
		gv, ok2 := variable.GetGlobalVarValue(name)
		if !ok2 {
			value, b := variable.GetEnvParamValue(name)
			return value, b
		} else {
			return gv, true
		}

	}
	return v, true
}

func (c *Context) SetVar(name string, value string) {
	//直接修改上下文的变量
	c.argsLock.Lock()
	defer c.argsLock.Unlock()
	if c.Args == nil {
		c.Args = map[string]string{}
	}
	c.Args[name] = value
}

func (c *Context) parseCode(code string) string{
	return variable.ExpandVars(code, func(name string)(v string){
		v, _=c.GetVar(name)
		return v
	})
}
