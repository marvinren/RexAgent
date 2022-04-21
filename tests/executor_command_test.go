package tests

import (
	"rexagent/pkg/executor/command"
	"rexagent/pkg/task"
	"strings"
	"testing"
)

func TestCommandExecutor(t *testing.T){
	job := &task.CommandJob {
		Id: "1",
		Name: "hello",
		Code: "echo ${app.name}",
		Lang: "exec",
		Timeout: 3,
		Parameters: []task.Parameter{
			{"app.name", "rexagent", ""},
		},
	}

	executor := &command.CommandExecutor{}

	ret, err := executor.Execute(job, nil)
	t.Log(string(ret))

	if err!=nil {
		t.Errorf("execute err: %s", err)
	}
	if strings.TrimSpace(string(ret)) != "rexagent" {
		t.Errorf("varibale is not corrrect")
	}
}
