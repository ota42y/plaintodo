package main

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

func TestExitCommand(t *testing.T) {
	cmd := NewExitCommand()

	cmds := make(map[string]Command)
	cmds["exit"] = cmd
	a := NewAutomaton(nil, cmds)

	terminate := a.Execute("exit")
	if !terminate {
		t.Errorf("ExitCommand.Execute shud be return true")
		t.FailNow()
	}
}

func TestReloadCommand(t *testing.T) {
	cmd := NewReloadCommand()

	cmds := make(map[string]Command)
	cmds["reload"] = cmd
	a := NewAutomaton(ReadTestConfig(), cmds)

	terminate := a.Execute("reload")
	if terminate {
		t.Errorf("ReloadCommand.Execute shud be return false")
		t.FailNow()
	}

	if len(a.Tasks) == 0 {
		t.Errorf("Task num shuldn't be 0")
		t.FailNow()
	}
}

func TestLsCommand(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := NewLsCommand(buf)

	cmds := make(map[string]Command)
	cmds["ls"] = cmd
	cmds["reload"] = NewReloadCommand()
	a := NewAutomaton(ReadTestConfig(), cmds)

	a.Execute("reload")
	terminate := a.Execute("ls")
	if terminate {
		t.Errorf("LsCommand.Execute shud be return false")
		t.FailNow()
	}

	if len(buf.String()) == 0 {
		t.Errorf("No outputs")
		t.FailNow()
	}
}

func TestLsAllCommand(t *testing.T) {
	cmds := make(map[string]Command)

	buf := &bytes.Buffer{}
	cmd := NewLsAllCommand(buf)
	cmds["lsall"] = cmd

	buf2 := &bytes.Buffer{}
	cmd2 := NewLsCommand(buf2)
	cmds["ls"] = cmd2

	cmds["reload"] = NewReloadCommand()
	a := NewAutomaton(ReadTestConfig(), cmds)

	a.Execute("reload")
	terminate := a.Execute("lsall")
	if terminate {
		t.Errorf("LsCommand.Execute shud be return false")
		t.FailNow()
	}

	a.Execute("ls")

	length := len(buf.String())

	if length == 0 {
		t.Errorf("No outputs")
		t.FailNow()
	}

	if length <= len(buf2.String()) {
		t.Errorf("lsall output %d length, but it's shuld be more longer than ls command output length (%d)", length, len(buf2.String()))
		t.FailNow()
	}
}

func TestCompleteCommand(t *testing.T) {
	cmds := make(map[string]Command)

	cmds["complete"] = NewCompleteCommand()

	cmds["reload"] = NewReloadCommand()
	a := NewAutomaton(ReadTestConfig(), cmds)
	a.Execute("reload")
	task := a.Tasks[0]

	if task.Attributes["complete"] != "" {
		t.Errorf("Task[\"complete\"} isn't blank")
		t.FailNow()
	}

	terminate := a.Execute(fmt.Sprintf("complete %d", task.Id))
	if terminate {
		t.Errorf("CompleteCommand.Execute shud be return false")
		t.FailNow()
	}

	_, err := time.Parse(dateTimeFormat, task.Attributes["complete"])
	if err != nil {
		t.Errorf("Task complete format invalid '%s'", task.Attributes["complete"])
		t.FailNow()
	}
}

func TestCompleteTask(t *testing.T) {
	tasks := ReadTestTasks()
	cmd := NewCompleteCommand()

	isComplete := cmd.completeTask(0, tasks)
	if isComplete {
		t.Errorf("If there is no task with taskId, completeTask shuld return false, but true")
		t.FailNow()
	}

	isComplete = cmd.completeTask(4, tasks)
	if !isComplete {
		t.Errorf("If there is task with taskId, completeTask shuld return true, but false")
		t.FailNow()
	}

	completeString := tasks[0].SubTasks[1].Attributes["complete"]
	_, err := time.Parse(dateTimeFormat, completeString)
	if err != nil {
		t.Errorf("Task complete format invalid '%s'", completeString)
		t.FailNow()
	}
}
