package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"./command"
	"./executor"
	"./task"
	"./util"
)

func TestGetIntAttribute(t *testing.T) {
	m := make(map[string]string)
	m["id"] = "42"

	n, err := util.GetIntAttribute("id", m)
	if err != nil {
		t.Errorf("When not error, shuld return err as nil, but %v", err)
		t.FailNow()
	}

	if 42 != n {
		t.Errorf("shuld return %d, but %d", 42, n)
		t.FailNow()
	}
}

func TestLsCommand(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := NewLsCommand(buf)

	config, _ := util.ReadTestConfig()

	cmds := make(map[string]command.Command)
	cmds["ls"] = cmd
	cmds["reload"] = command.NewReload()
	a := executor.NewExecutor(config, cmds)

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
	cmds := make(map[string]command.Command)

	buf := &bytes.Buffer{}
	cmd := NewLsAllCommand(buf)
	cmds["lsall"] = cmd

	buf2 := &bytes.Buffer{}
	cmd2 := NewLsCommand(buf2)
	cmds["ls"] = cmd2

	cmds["reload"] = command.NewReload()

	config, _ := util.ReadTestConfig()
	a := executor.NewExecutor(config, cmds)

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

func TestAddSubTaskCommand(t *testing.T) {
	taskName := "create sub task"
	taskStart := "2015-02-01"

	cmds := make(map[string]command.Command)
	cmds["task"] = command.NewAddTask()
	cmds["subtask"] = NewAddSubTaskCommand()
	cmds["reload"] = command.NewReload()

	config, buf := util.ReadTestConfig()
	a := executor.NewExecutor(config, cmds)
	a.Execute("reload")
	buf.Reset()

	taskID := a.S.MaxTaskID

	input := "subtask 6 " + taskName + " :start " + taskStart
	terminate := a.Execute(input)

	if terminate {
		t.Errorf("AddSubTaskCommand terminate automaton")
		t.FailNow()
	}

	parent := a.S.Tasks[0].SubTasks[1].SubTasks[1]

	if len(parent.SubTasks) == 0 {
		t.Errorf("SubTask not add")
		t.FailNow()
	}

	task := parent.SubTasks[0]
	if task.Level != parent.Level+1 {
		t.Errorf("Subtask level shuld be %d but %d", parent.Level+1, task.Level)
		t.FailNow()
	}

	if task.Name != taskName {
		t.Errorf("Task name shud %s, but %s", taskName, task.Name)
		t.FailNow()
	}

	if task.Attributes["start"] != taskStart {
		t.Errorf("Task start shud %s, but %s", taskStart, task.Attributes["start"])
		t.FailNow()
	}

	if taskID+1 != task.ID {
		t.Errorf("Task's id shud be %d but %d", taskID+1, task.ID)
		t.FailNow()
	}

	if a.S.MaxTaskID != task.ID {
		t.Errorf("Automaton.MaxTaskID shuld be %d, but %d", a.S.MaxTaskID, task.ID)
		t.FailNow()
	}
	taskID = a.S.MaxTaskID

	outputString := buf.String()
	correctString := "subtask hit\nCreate SubTask:\nParent: " + parent.String(true) + "\nSubTask: " + task.String(true) + "\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	buf = &bytes.Buffer{}
	config.Writer = buf
	a.Execute("subtask test")

	outputString = buf.String()
	correctString = "subtask hit\nCreate Subtask error: invalid format 'test'\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	buf = &bytes.Buffer{}
	config.Writer = buf
	a.Execute("subtask 0 no parent task")

	outputString = buf.String()
	correctString = "subtask hit\nCreate SubTask error: thee is no task which have :id 0\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	if a.S.MaxTaskID != taskID {
		t.Errorf("When error occerd, Automaton.MaxTaskID shuldn't change but %d", taskID)
		t.FailNow()
	}

	a.Execute("task test11")
	a.Execute("task test12")
	a.Execute("subtask 12 child in task12")
	if len(a.S.Tasks[3].SubTasks) == 0 {
		t.Errorf("Subtask isn't added by id=12 task's child")
		t.FailNow()
	}
}

func TestStartCommand(t *testing.T) {
	cmd := NewStartCommand()

	cmds := make(map[string]command.Command)
	cmds["reload"] = command.NewReload()
	cmds["start"] = cmd
	config, _ := util.ReadTestConfig()
	a := executor.NewExecutor(config, cmds)
	a.Execute("reload")

	now := time.Now()

	task := a.S.Tasks[0].SubTasks[1].SubTasks[0]
	if _, ok := task.Attributes["start"]; ok {
		t.Errorf("task already set start attribute, test data is invalid %v", task)
		t.FailNow()
	}

	terminate := a.Execute("start :id 5")
	if terminate {
		t.Errorf("StartCommand.Execute shud be return false")
		t.FailNow()
	}

	value, ok := task.Attributes["start"]
	if !ok {
		t.Errorf("start attribute not set")
		t.FailNow()
	}

	dateTime, ok := util.ParseTime(value)
	diff := dateTime.Sub(now)
	if diff.Minutes() < -2 || 2 < diff.Minutes() {
		t.Errorf("set time (%v) isn't now because %v minutes after", value, diff.Seconds())
		t.FailNow()
	}

	task.Attributes["start"] = time.Now().AddDate(1, 0, 0).Format(util.DateTimeFormat)
	terminate = a.Execute("start :id 5")
	dateTime, ok = util.ParseTime(value)
	diff = dateTime.Sub(now)
	if diff.Minutes() < -2 || 2 < diff.Minutes() {
		t.Errorf("set new start time, but old time isn't overwrited")
		t.FailNow()
	}
}

func TestMoveCommand(t *testing.T) {
	cmd := NewMoveCommand()

	cmds := make(map[string]command.Command)
	cmds["reload"] = command.NewReload()
	cmds["move"] = cmd
	config, buf := util.ReadTestConfig()
	a := executor.NewExecutor(config, cmds)
	a.Execute("reload")
	buf.Reset()

	fromTask, moveTask := task.GetTask(4, a.S.Tasks)
	_, toTask := task.GetTask(8, a.S.Tasks)

	fromNum := len(fromTask.SubTasks)
	toNum := len(toTask.SubTasks)

	a.Execute("move :to 42")
	outputString := buf.String()
	correctString := fmt.Sprintf("move hit\nnot exist :from\n")
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	a.Execute("move :from 4")
	outputString = buf.String()
	correctString = fmt.Sprintf("move hit\nnot exist :to\n")
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	a.Execute("move :from 42 :to 4")
	outputString = buf.String()
	correctString = fmt.Sprintf("move hit\nthere is no exist %d task\n", 42)
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	a.Execute("move :from 4 :to 42")
	outputString = buf.String()
	correctString = fmt.Sprintf("move hit\nthere is no exist %d task\n", 42)
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	movedParent, movedTask := task.GetTask(4, a.S.Tasks)
	if movedParent == nil {
		t.Errorf("when not meved task, parent shuldn't be change from %v but %v", toTask, movedParent)
		t.FailNow()
	}
	if movedParent.ID != fromTask.ID {
		t.Errorf("when not meved task, parent shuldn't be change from %v but %v", fromTask, movedParent)
		t.FailNow()
	}

	terminate := a.Execute("move :from 4 :to 8")
	if terminate {
		t.Errorf("ReloadCommand.Execute shud be return false")
		t.FailNow()
	}
	outputString = buf.String()
	correctString = fmt.Sprintf("move hit\ntask moved to sub task\nparent: %v\n", toTask.String(true))
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	if fromNum-1 != len(fromTask.SubTasks) {
		t.Errorf("if task moved, from task's subtask shuld be %d, but %d in %v", fromNum-1, len(fromTask.SubTasks), fromTask)
		t.FailNow()
	}

	if toNum+1 != len(toTask.SubTasks) {
		t.Errorf("if task moved, to task's subtask shuld be %d, but %d", toNum+1, len(toTask.SubTasks))
		t.FailNow()
	}

	movedParent, movedTask = task.GetTask(4, a.S.Tasks)
	if toTask.ID != movedParent.ID {
		t.Errorf("move %d task's sub task, but %d task's subtask", toTask.ID, movedParent.ID)
		t.FailNow()
	}

	if 3 != len(movedTask.SubTasks) {
		t.Errorf("move all sub tasks (%d), but %d sub task exist", 3, len(movedTask.SubTasks))
		t.FailNow()
	}

	a.Execute("move :from 4 :to 9")
	if 2 != moveTask.Level {
		t.Errorf("when task moved, Task.Level shuld be %d, but %d", 2, moveTask.Level)
		t.FailNow()
	}
	if 3 != moveTask.SubTasks[0].Level {
		t.Errorf("when task moved, sub task's Task.Level shuld be %d, but %d", 3, moveTask.SubTasks[0].Level)
		t.FailNow()
	}

	buf.Reset()

	a.Execute("move :from 4 :to 0")
	outputString = buf.String()
	correctString = fmt.Sprintf("move hit\ntask moved to top level task\n")
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	movedParent, movedTask = task.GetTask(4, a.S.Tasks)
	if movedParent != nil {
		t.Errorf("if task moved top level task, parent shuld be nil but %v", movedParent)
		t.FailNow()
	}
	if len(a.S.Tasks) != 3 {
		t.Errorf("if task moved top level task, Automaton.Task shuld be %d tasks, but %d", 3, len(a.S.Tasks))
		t.FailNow()
	}

	if 0 != moveTask.Level {
		t.Errorf("when task moved, Task.Level shuld be %d, but %d", 0, moveTask.Level)
		t.FailNow()
	}

	if 1 != moveTask.SubTasks[0].Level {
		t.Errorf("when task moved, sub task's Task.Level shuld be %d, but %d", 1, moveTask.SubTasks[0].Level)
		t.FailNow()
	}
}

func TestOpenCommand(t *testing.T) {
	cmd := NewOpenCommand()

	cmds := make(map[string]command.Command)
	cmds["reload"] = command.NewReload()
	cmds["open"] = cmd
	config, buf := util.ReadTestConfig()
	a := executor.NewExecutor(config, cmds)
	a.Execute("reload")
	buf.Reset()

	a.Execute("open :id 8")
	_, tk := task.GetTask(8, a.S.Tasks)
	outputString := buf.String()
	correctString := fmt.Sprintf("open hit\nThere is no url in task:\n%s\n", tk.String(true))
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	a.Execute("open :id 42")
	outputString = buf.String()
	correctString = fmt.Sprintf("open hit\nThere is no such task :id %d\n", 42)
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	a.Execute("open :aaa bbb")
	outputString = buf.String()
	correctString = fmt.Sprintf("open hit\nnot exist :id\n")
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	_, tk = task.GetTask(9, a.S.Tasks)
	url, err := cmd.getUrl(tk)
	if tk.Attributes["url"] != url {
		t.Errorf("getUrl shuld return %s, but %s", tk.Attributes["url"], url)
		t.FailNow()
	}

	if err != nil {
		t.Errorf("getUrl shuldn't return erro but %s", err.Error())
		t.FailNow()
	}

	tk = a.S.Tasks[0]
	url, err = cmd.getUrl(tk)
	outputString = err.Error()
	correctString = fmt.Sprintf("There is no url in task:\n%s", tk.String(true))
	if correctString != outputString {
		t.Errorf("if task haven't :url attribute and task name is url, getUrl shuld return error %s, but %s", correctString, outputString)
		t.FailNow()
	}
}

type CommandTest struct {
	T         *testing.T
	Option    string
	Called    bool
	Terminate bool
}

func (t *CommandTest) Execute(option string, s *command.State) (terminate bool) {
	t.Called = true

	if option != t.Option {
		t.T.Errorf("option shud be %s but %s", t.Option, option)
		t.T.FailNow()
	}

	return t.Terminate
}

func TestAliasCommand(t *testing.T) {
	cmds := make(map[string]command.Command)
	cmds["alias"] = NewAliasCommand()
	config, buf := util.ReadTestConfig()
	a := executor.NewExecutor(config, cmds)

	terminate := a.Execute("alias")

	outputString := buf.String()
	correctString := fmt.Sprintf("alias hit\nlsalltest = ls :level 1\npt = postpone :postpone 1 day :id\n")
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	if terminate {
		t.Errorf("shud be return false")
		t.FailNow()
	}

	cmd := &CommandTest{
		T:         t,
		Option:    ":postpone 1 day :id 1",
		Called:    false,
		Terminate: false,
	}

	cmds["postpone"] = cmd
	// postpone tomorrow
	terminate = a.Execute("pt 1")

	if terminate != cmd.Terminate {
		t.Errorf("Automation.Execute shud be return %v but %v", terminate, cmd.Terminate)
		t.FailNow()
	}

	if !cmd.Called {
		t.Errorf("command not called")
		t.FailNow()
	}

	outputString = buf.String()
	correctString = fmt.Sprintf("alias pt = postpone :postpone 1 day :id\ncommand: postpone :postpone 1 day :id 1\npostpone hit")
	if !strings.HasPrefix(outputString, correctString) {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()
}
