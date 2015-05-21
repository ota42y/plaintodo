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

func TestCompleteCommandError(t *testing.T) {
	cmds := make(map[string]Command)

	cmds["complete"] = NewCompleteCommand()

	cmds["reload"] = NewReloadCommand()
	config := ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")

	buf := &bytes.Buffer{}
	config.Writer = buf

	terminate := a.Execute(fmt.Sprintf("complete %da", 1))
	if terminate {
		t.Errorf("CompleteCommand.Execute shud be return false")
		t.FailNow()
	}

	outputString := buf.String()
	if outputString != "complete hit\nstrconv.ParseInt: parsing \"1a\": invalid syntax" {
		t.Errorf("CompleteCommand.Execute shuld write error, but %s", outputString)
		t.FailNow()
	}

	buf.Reset()
	terminate = a.Execute(fmt.Sprintf("complete %d", 100))
	if terminate {
		t.Errorf("CompleteCommand.Execute shud be return false")
		t.FailNow()
	}

	outputString = buf.String()
	if outputString != "complete hit\nThere is no Task which have task id: 100\n" {
		t.Errorf("CompleteCommand.Execute shuld write no such task error, but %s", outputString)
		t.FailNow()
	}
}

func TestCompleteCommand(t *testing.T) {
	cmds := make(map[string]Command)

	cmds["complete"] = NewCompleteCommand()

	cmds["reload"] = NewReloadCommand()
	config := ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")
	task := a.Tasks[0]

	if task.Attributes["complete"] != "" {
		t.Errorf("Task[\"complete\"} isn't blank")
		t.FailNow()
	}

	buf := &bytes.Buffer{}
	config.Writer = buf

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

	outputString := buf.String()
	correctString := fmt.Sprintf("complete hit\nComplete %s and %d sub tasks\n", task.Name, 7)
	if outputString != correctString {
		t.Errorf("CompleteCommand.Execute shuld write '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
}

func isAllCompleted(task *Task) bool {
	_, ok := task.Attributes["complete"]
	if !ok {
		return false
	}

	for _, subTask := range task.SubTasks {
		if !isAllCompleted(subTask) {
			return false
		}
	}

	return true
}

func TestCompleteTask(t *testing.T) {
	tasks := ReadTestTasks()
	cmd := NewCompleteCommand()

	completeTask, n := cmd.completeTask(0, tasks)
	if completeTask != nil {
		t.Errorf("If there is no task with taskId, completeTask shuld return nil, but %v", completeTask)
		t.FailNow()
	}

	if n != 0 {
		t.Errorf("If there is no task with taskId, completeTask shuld return complete 0 subtask, but %d", n)
		t.FailNow()
	}

	alreadyCompleted := "2014-01-01"
	tasks[0].SubTasks[1].SubTasks[1].Attributes["complete"] = alreadyCompleted

	completeTask, n = cmd.completeTask(4, tasks)
	if completeTask == nil {
		t.Errorf("If there is task with taskId, completeTask shuld return complete task, but nil")
		t.FailNow()
	}

	if n != 3 {
		t.Errorf("If there is task with taskId, completeTask shuld return complete subtask num (4) but %d", n)
		t.FailNow()
	}

	if !isAllCompleted(tasks[0].SubTasks[1]) {
		t.Errorf("not complete selected task and all sub tasks")
		t.FailNow()
	}

	completeString := tasks[0].SubTasks[1].Attributes["complete"]
	_, err := time.Parse(dateTimeFormat, completeString)
	if err != nil {
		t.Errorf("Task complete format invalid '%s'", completeString)
		t.FailNow()
	}

	completeString = tasks[0].SubTasks[1].SubTasks[0].Attributes["complete"]
	_, err = time.Parse(dateTimeFormat, completeString)
	if err != nil {
		t.Errorf("Task complete format invalid '%s'", completeString)
		t.FailNow()
	}

	alreadyCompletedAttribute := tasks[0].SubTasks[1].SubTasks[1].Attributes["complete"]
	if (alreadyCompleted != alreadyCompletedAttribute) || alreadyCompletedAttribute == completeString {
		t.Errorf("Already completed task isn't overwrite but %s", alreadyCompletedAttribute)
		t.FailNow()
	}
}

func TestGetCompleteDayList(t *testing.T) {
	tasks := ReadTestTasks()
	cmd := NewSaveCommand()

	testTimeList := [...]string{"2015-01-31 10:42", "2015-01-29", "2015-01-30 10:42", "2015-01-30"}
	tasks[0].Attributes["complete"] = testTimeList[0]
	tasks[0].SubTasks[0].Attributes["complete"] = testTimeList[1]
	tasks[0].SubTasks[1].Attributes["complete"] = testTimeList[2]
	tasks[0].SubTasks[1].SubTasks[0].Attributes["complete"] = testTimeList[3]

	correctTimeList := make([]time.Time, 3)
	parseList := [...]int{1, 2, 0}
	for index, value := range parseList {
		timeData, ok := ParseTime(testTimeList[value])
		if !ok {
			t.Errorf("parse error %s", testTimeList[value])
			t.FailNow()
		}
		correctTimeList[index] = timeData
	}

	timeList := cmd.getCompleteDayList(tasks)
	if len(timeList) != len(correctTimeList) {
		t.Errorf("shuld return %d items, but %d items %v", len(correctTimeList), len(timeList), timeList)
		t.FailNow()
	}

	for index, item := range correctTimeList {
		year, month, day := item.Date()
		y, m, d := timeList[index].Date()

		if (year != y) || (month != m) || (day != d) {
			t.Errorf("shuld return %v but %v", item, timeList[index])
			t.FailNow()
		}
	}
}

