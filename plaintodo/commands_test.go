package main

import (
	"bytes"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"

	"./task"
	"./util"
)

func TestGetIntAttribute(t *testing.T) {
	m := make(map[string]string)
	m["id"] = "42"

	n, err := GetIntAttribute("id", m)
	if err != nil {
		t.Errorf("When not error, shuld return err as nil, but %v", err)
		t.FailNow()
	}

	if 42 != n {
		t.Errorf("shuld return %d, but %d", 42, n)
		t.FailNow()
	}
}

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
	a := NewAutomaton(util.ReadTestConfig(), cmds)

	terminate := a.Execute("reload")
	if terminate {
		t.Errorf("ReloadCommand.Execute shud be return false")
		t.FailNow()
	}

	if len(a.Tasks) == 0 {
		t.Errorf("Task num shuldn't be 0")
		t.FailNow()
	}

	id := a.Tasks[1].SubTasks[0].ID
	if a.MaxTaskID != id {
		t.Errorf("Save max task id %d, but %d", id, a.MaxTaskID)
		t.FailNow()
	}
}

func TestLsCommand(t *testing.T) {
	buf := &bytes.Buffer{}
	cmd := NewLsCommand(buf)

	cmds := make(map[string]Command)
	cmds["ls"] = cmd
	cmds["reload"] = NewReloadCommand()
	a := NewAutomaton(util.ReadTestConfig(), cmds)

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
	a := NewAutomaton(util.ReadTestConfig(), cmds)

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
	config := util.ReadTestConfig()
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
	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")
	task := a.Tasks[0]

	if task.Attributes["complete"] != "" {
		t.Errorf("Task[\"complete\"} isn't blank")
		t.FailNow()
	}

	buf := &bytes.Buffer{}
	config.Writer = buf

	terminate := a.Execute(fmt.Sprintf("complete %d", task.ID))
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

func isAllCompleted(t *task.Task) bool {
	_, ok := t.Attributes["complete"]
	if !ok {
		return false
	}

	for _, subTask := range t.SubTasks {
		if !isAllCompleted(subTask) {
			return false
		}
	}

	return true
}

func TestSetNewRepeat(t *testing.T) {
	cmd := NewCompleteCommand()
	tk := &task.Task{
		Attributes: make(map[string]string),
	}

	now := time.Now()
	base := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local)
	tk.Attributes["start"] = base.Format(dateTimeFormat)

	tk.Attributes["repeat"] = "every 1 day"
	cmd.setNewRepeat(now, tk)
	correct := base.AddDate(0, 0, 1)
	correctString := correct.Format(dateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(dateTimeFormat)
	tk.Attributes["repeat"] = "every 1 month"
	cmd.setNewRepeat(now, tk)
	correct = base.AddDate(0, 1, 0)
	correctString = correct.Format(dateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(dateTimeFormat)
	tk.Attributes["repeat"] = "every 1 year"
	cmd.setNewRepeat(now, tk)
	correct = base.AddDate(1, 0, 0)
	correctString = correct.Format(dateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(dateTimeFormat)
	tk.Attributes["repeat"] = "every 2 week"
	cmd.setNewRepeat(now, tk)
	correct = base.AddDate(0, 0, 14)
	correctString = correct.Format(dateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(dateTimeFormat)
	tk.Attributes["repeat"] = "every 30 minutes"
	cmd.setNewRepeat(now, tk)
	correct = base.Add(30 * time.Minute)
	correctString = correct.Format(dateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(dateTimeFormat)
	tk.Attributes["repeat"] = "every 2 hour"
	cmd.setNewRepeat(now, tk)
	correct = base.Add(2 * time.Hour)
	correctString = correct.Format(dateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(dateTimeFormat)
	tk.Attributes["repeat"] = "after 4 day"
	cmd.setNewRepeat(now, tk)
	correct = now.AddDate(0, 0, 4)
	correctString = correct.Format(dateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}
}

func TestCompleteRepeatTask(t *testing.T) {
	tasks := util.ReadTestTasks()
	baseTask := tasks[1].SubTasks[0]

	postponeCommand := NewPostponeCommand()
	optionMap := make(map[string]string)
	optionMap["postpone"] = "1 day"
	postponeCommand.postpone(baseTask, optionMap)

	cmd := NewCompleteCommand()
	completeTask, newTasks, n := cmd.completeTask(8, tasks)
	if completeTask == nil {
		t.Errorf("If there is task with taskID, completeTask shuld return complete task, but nil")
		t.FailNow()
	}

	if n != 2 {
		t.Errorf("If there is task with taskID, completeTask shuld return complete subtask num (2) but %d", n)
		t.FailNow()
	}

	if len(newTasks) != 3 {
		t.Errorf("If repeat task complete, task will copy")
		t.FailNow()
	}

	baseStart, baseOk := ParseTime(newTasks[1].SubTasks[0].Attributes["start"])
	repeatStart, repeatOk := ParseTime(newTasks[2].SubTasks[0].Attributes["start"])
	if !baseOk || !repeatOk {
		t.Errorf("start parse error")
		t.FailNow()
	}

	nextStart := baseStart.AddDate(0, 0, 1)
	if nextStart != repeatStart {
		t.Errorf("set after 1 day (%v), but %v", nextStart, repeatStart)
		t.FailNow()
	}

	if _, ok := newTasks[1].SubTasks[0].Attributes["postpone"]; !ok {
		t.Errorf("postpone attribute delete from base task %v", newTasks[1])
		t.FailNow()
	}

	if _, ok := newTasks[2].SubTasks[0].Attributes["postpone"]; ok {
		t.Errorf("postpone attribute copy to repeat task")
		t.FailNow()
	}
	delete(newTasks[1].SubTasks[0].Attributes, "postpone")

	delete(newTasks[1].SubTasks[0].Attributes, "start")
	delete(newTasks[2].SubTasks[0].Attributes, "start")
	delete(newTasks[1].Attributes, "complete")
	delete(newTasks[1].SubTasks[0].Attributes, "complete")
	if !newTasks[1].Equal(newTasks[2]) {
		t.Errorf("If copy by repeat, it's same task without complete attribute")
		t.FailNow()
	}
}

func TestCompleteTask(t *testing.T) {
	tasks := util.ReadTestTasks()
	cmd := NewCompleteCommand()

	completeTask, tasks, n := cmd.completeTask(0, tasks)
	if completeTask != nil {
		t.Errorf("If there is no task with taskID, completeTask shuld return nil, but %v", completeTask)
		t.FailNow()
	}

	if len(tasks) != 2 {
		t.Errorf("task num shudn't change")
		t.FailNow()
	}

	if n != 0 {
		t.Errorf("If there is no task with taskID, completeTask shuld return complete 0 subtask, but %d", n)
		t.FailNow()
	}

	alreadyCompleted := "2014-01-01"
	tasks[0].SubTasks[1].SubTasks[1].Attributes["complete"] = alreadyCompleted

	completeTask, tasks, n = cmd.completeTask(4, tasks)
	if completeTask == nil {
		t.Errorf("If there is task with taskID, completeTask shuld return complete task, but nil")
		t.FailNow()
	}

	if len(tasks) != 2 {
		t.Errorf("task num shudn't change")
		t.FailNow()
	}

	if n != 3 {
		t.Errorf("If there is task with taskID, completeTask shuld return complete subtask num (4) but %d", n)
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
	tasks := util.ReadTestTasks()
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

func TestAddTaskCommand(t *testing.T) {
	taskName := "create new task"
	taskStart := "2015-02-01"

	cmds := make(map[string]Command)
	cmds["task"] = NewAddTaskCommand()
	cmds["reload"] = NewReloadCommand()

	config := util.ReadTestConfig()
	buf := &bytes.Buffer{}
	config.Writer = buf

	a := NewAutomaton(config, cmds)

	input := "task " + taskName + " :start " + taskStart
	terminate := a.Execute(input)

	if terminate {
		t.Errorf("AddTaskCommand terminate automaton")
		t.FailNow()
	}

	if len(a.Tasks) == 0 {
		t.Errorf("Task not add")
		t.FailNow()
	}

	task := a.Tasks[0]
	if task.Name != taskName {
		t.Errorf("Task name shud %s, but %s", taskName, task.Name)
		t.FailNow()
	}

	if task.Attributes["start"] != taskStart {
		t.Errorf("Task start shud %s, but %s", taskStart, task.Attributes["start"])
		t.FailNow()
	}

	if a.MaxTaskID != task.ID {
		t.Errorf("Automaton.MaxTaskID shuld be %d, but %d", a.MaxTaskID, task.ID)
		t.FailNow()
	}

	outputString := buf.String()
	correctString := "task hit\nCreate task: " + taskName + " :id 1 :start " + taskStart + "\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	taskID := a.MaxTaskID

	buf = &bytes.Buffer{}
	config.Writer = buf
	a.Execute("task ")

	outputString = buf.String()
	correctString = "task hit\nCreate task error: blank line\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	if a.MaxTaskID != taskID {
		t.Errorf("When error occerd, Automaton.MaxTaskID shuldn't change but %d", taskID)
		t.FailNow()
	}
}

func TestAddSubTaskCommand(t *testing.T) {
	taskName := "create sub task"
	taskStart := "2015-02-01"

	cmds := make(map[string]Command)
	cmds["task"] = NewAddTaskCommand()
	cmds["subtask"] = NewAddSubTaskCommand()
	cmds["reload"] = NewReloadCommand()

	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")

	buf := &bytes.Buffer{}
	config.Writer = buf

	taskID := a.MaxTaskID

	input := "subtask 6 " + taskName + " :start " + taskStart
	terminate := a.Execute(input)

	if terminate {
		t.Errorf("AddSubTaskCommand terminate automaton")
		t.FailNow()
	}

	parent := a.Tasks[0].SubTasks[1].SubTasks[1]

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

	if a.MaxTaskID != task.ID {
		t.Errorf("Automaton.MaxTaskID shuld be %d, but %d", a.MaxTaskID, task.ID)
		t.FailNow()
	}
	taskID = a.MaxTaskID

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

	if a.MaxTaskID != taskID {
		t.Errorf("When error occerd, Automaton.MaxTaskID shuldn't change but %d", taskID)
		t.FailNow()
	}

	a.Execute("task test11")
	a.Execute("task test12")
	a.Execute("subtask 12 child in task12")
	if len(a.Tasks[3].SubTasks) == 0 {
		t.Errorf("Subtask isn't added by id=12 task's child")
		t.FailNow()
	}
}

func TestSetAttributeCommand(t *testing.T) {
	cmd := NewSetAttributeCommand()
	url := "http://example.com"

	cmds := make(map[string]Command)
	cmds["reload"] = NewReloadCommand()
	cmds["set"] = cmd
	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")

	buf := &bytes.Buffer{}
	config.Writer = buf

	terminate := a.Execute("set :url " + url)
	if terminate {
		t.Errorf("SetAttributeCommand.Execute shud be return false")
		t.FailNow()
	}

	outputString := buf.String()
	correctString := "set hit\nnot exist :id\n"
	if outputString != correctString {
		t.Errorf("Shuld output '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}

	task := a.Tasks[0].SubTasks[0].SubTasks[0]
	task.Attributes["url"] = url
	correctString = "set hit\nset attribute " + task.String(true) + "\n"
	delete(task.Attributes, "url")

	buf.Reset()
	terminate = a.Execute("set :id 3 :url " + url)
	if terminate {
		t.Errorf("SetAttributeCommand.Execute shud be return false")
		t.FailNow()
	}

	value, ok := task.Attributes["url"]
	if !ok {
		t.Errorf("attribute not set")
		t.FailNow()
	}

	if value != url {
		t.Errorf("set attribute shuld %s, but %s", url, value)
		t.FailNow()
	}

	outputString = buf.String()
	if outputString != correctString {
		t.Errorf("Shuld output '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}

	buf.Reset()
	terminate = a.Execute("set :id 0 :url " + url)

	outputString = buf.String()
	correctString = "set hit\nthere is no exist :id 0 task\n"
	if outputString != correctString {
		t.Errorf("Shuld output '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
}

func TestStartCommand(t *testing.T) {
	cmd := NewStartCommand()

	cmds := make(map[string]Command)
	cmds["reload"] = NewReloadCommand()
	cmds["start"] = cmd
	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")

	buf := &bytes.Buffer{}
	config.Writer = buf
	now := time.Now()

	task := a.Tasks[0].SubTasks[1].SubTasks[0]
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

	dateTime, ok := ParseTime(value)
	diff := dateTime.Sub(now)
	if diff.Minutes() < -2 || 2 < diff.Minutes() {
		t.Errorf("set time (%v) isn't now because %v minutes after", value, diff.Seconds())
		t.FailNow()
	}

	task.Attributes["start"] = time.Now().AddDate(1, 0, 0).Format(dateTimeFormat)
	terminate = a.Execute("start :id 5")
	dateTime, ok = ParseTime(value)
	diff = dateTime.Sub(now)
	if diff.Minutes() < -2 || 2 < diff.Minutes() {
		t.Errorf("set new start time, but old time isn't overwrited")
		t.FailNow()
	}
}

func TestPostponeCommand(t *testing.T) {
	cmd := NewPostponeCommand()

	cmds := make(map[string]Command)
	cmds["reload"] = NewReloadCommand()
	cmds["start"] = NewStartCommand()
	cmds["postpone"] = cmd
	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")

	buf := &bytes.Buffer{}
	config.Writer = buf
	now := time.Now()

	tk := a.Tasks[0].SubTasks[1].SubTasks[0]
	if _, ok := tk.Attributes["postpone"]; ok {
		t.Errorf("task already set postpone attribute, test data is invalid %v", tk)
		t.FailNow()
	}

	a.Execute("postpone :id 5 :postpone 1 month")
	outputString := buf.String()
	correctString := fmt.Sprintln("postpone hit\ntask :id", tk.ID, "haven't start attribute, so postpone not work")
	if outputString != correctString {
		t.Errorf("shuld return '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	tk.Attributes["start"] = "test"
	a.Execute("postpone :id 5 :postpone 1 month")
	outputString = buf.String()
	correctString = fmt.Sprintln("postpone hit\ntest is invalid format, so postpone not work")
	if outputString != correctString {
		t.Errorf("shuld return '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	// set start 0 time
	tk.Attributes["start"] = time.Unix(0, 0).Format(dateTimeFormat)

	// invalid case
	a.Execute("postpone :id 5 :postpone 1")
	outputString = buf.String()
	correctString = fmt.Sprintln("postpone hit\n1 is invalid format")
	if outputString != correctString {
		t.Errorf("shuld return '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}

	terminate := a.Execute("postpone :id 5 :postpone 1 month")
	if terminate {
		t.Errorf("PostPoneCommand.Execute shud be return false")
		t.FailNow()
	}

	value, ok := tk.Attributes["postpone"]
	if !ok {
		t.Errorf("postpone attribute not set %v", tk)
		t.FailNow()
	}

	dateTime, ok := ParseTime(value)
	if !ok {
		t.Errorf("postpone attribute value is invalid formt %s", value)
		t.FailNow()
	}

	diff := dateTime.Sub(now.AddDate(0, 1, 0))
	if diff.Minutes() < -2 || 2 < diff.Minutes() {
		t.Errorf("postpone time (%v) isn't 1 month ofter because %v minutes after", value, diff.Minutes())
		t.FailNow()
	}
}

func TestMoveCommand(t *testing.T) {
	cmd := NewMoveCommand()

	cmds := make(map[string]Command)
	cmds["reload"] = NewReloadCommand()
	cmds["move"] = cmd
	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")

	buf := &bytes.Buffer{}
	config.Writer = buf

	fromTask, moveTask := task.GetTask(4, a.Tasks)
	_, toTask := task.GetTask(8, a.Tasks)

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

	movedParent, movedTask := task.GetTask(4, a.Tasks)
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

	movedParent, movedTask = task.GetTask(4, a.Tasks)
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

	movedParent, movedTask = task.GetTask(4, a.Tasks)
	if movedParent != nil {
		t.Errorf("if task moved top level task, parent shuld be nil but %v", movedParent)
		t.FailNow()
	}
	if len(a.Tasks) != 3 {
		t.Errorf("if task moved top level task, Automaton.Task shuld be %d tasks, but %d", 3, len(a.Tasks))
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

	cmds := make(map[string]Command)
	cmds["reload"] = NewReloadCommand()
	cmds["open"] = cmd
	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")

	buf := &bytes.Buffer{}
	config.Writer = buf

	a.Execute("open :id 8")
	_, tk := task.GetTask(8, a.Tasks)
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

	_, tk = task.GetTask(9, a.Tasks)
	url, err := cmd.getUrl(tk)
	if tk.Attributes["url"] != url {
		t.Errorf("getUrl shuld return %s, but %s", tk.Attributes["url"], url)
		t.FailNow()
	}

	if err != nil {
		t.Errorf("getUrl shuldn't return erro but %s", err.Error())
		t.FailNow()
	}

	tk = a.Tasks[0]
	url, err = cmd.getUrl(tk)
	outputString = err.Error()
	correctString = fmt.Sprintf("There is no url in task:\n%s", tk.String(true))
	if correctString != outputString {
		t.Errorf("if task haven't :url attribute and task name is url, getUrl shuld return error %s, but %s", correctString, outputString)
		t.FailNow()
	}
}

func TestNiceCommand(t *testing.T) {
	cmd := NewNiceCommand()

	cmds := make(map[string]Command)
	cmds["reload"] = NewReloadCommand()
	cmds["task"] = NewAddTaskCommand()
	cmds["nice"] = cmd
	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)
	a.Execute("reload")

	buf := &bytes.Buffer{}
	config.Writer = buf

	evernoteUrl := "https://www.evernote.com/shard/s1/nl/111111/abfdef-ght1234567890"
	a.Execute("task test :url " + evernoteUrl)
	buf.Reset()

	_, task := task.GetTask(a.MaxTaskID, a.Tasks)
	if evernoteUrl != task.Attributes["url"] {
		t.Errorf(":url attribute shuld be %s, but %s, test data invalid", evernoteUrl, task.Attributes["url"])
		t.FailNow()
	}

	terminate := a.Execute(fmt.Sprintf("nice :id %d", task.ID))

	outputString := buf.String()
	correctString := fmt.Sprintf("nice hit\nDone nice\nevernote url change 1 tasks\n")
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	if terminate {
		t.Errorf("shud be return false")
		t.FailNow()
	}

	if !strings.HasPrefix(task.Attributes["url"], "evernote:///view/") {
		t.Errorf(":url shuld have prefix '%s', but '%s'", "evernote:///view/", task.Attributes["url"])
		t.FailNow()
	}

	// if there is / in bottom
	task.Attributes["url"] = evernoteUrl + "/"
	// if sub task
	a.Tasks[0].SubTasks[0].Attributes["url"] = evernoteUrl

	a.Execute(fmt.Sprintf("nice"))
	outputString = buf.String()
	correctString = fmt.Sprintf("nice hit\nDone nice\nevernote url change 2 tasks\n")
	if outputString != correctString {
		t.Errorf("output shuld be '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	if !strings.HasPrefix(task.Attributes["url"], "evernote:///view/") {
		t.Errorf(":url shuld have prefix '%s', but '%s'", "evernote:///view/", task.Attributes["url"])
		t.FailNow()
	}

	if !strings.HasPrefix(a.Tasks[0].SubTasks[0].Attributes["url"], "evernote:///view/") {
		t.Errorf(":url shuld have prefix '%s', but '%s'", "evernote:///view/", a.Tasks[0].Attributes["url"])
		t.FailNow()
	}

	if a.Tasks[1].SubTasks[0].Attributes["url"] != "http://ota42y.com" {
		t.Errorf("if :url isn't evernote url, nice command not chenga, but change to %s", a.Tasks[1].SubTasks[0].Attributes["url"])
		t.FailNow()
	}
}

func TestAliasCommand(t *testing.T) {
	cmds := make(map[string]Command)
	cmds["alias"] = NewAliasCommand()
	config := util.ReadTestConfig()
	a := NewAutomaton(config, cmds)

	buf := &bytes.Buffer{}
	config.Writer = buf

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

func TestSaveCommand(t *testing.T) {
	Convey("correct", t, func() {
		Convey("save tasks", func() {
			cmd := NewSaveCommand()
			emptyTasks := make([]*task.Task, 0)
			cmd.saveToFile(emptyTasks, "../result/savetest.txt")

			tasks := util.ReadTestTasks()
			tasks[0].SubTasks = emptyTasks
			cmd.saveToFile(tasks, "../result/savetest.txt")

			loadTasks, _ := task.ReadTasks("../result/savetest.txt")
			So(len(loadTasks), ShouldEqual, len(tasks))
			So(len(loadTasks[0].SubTasks), ShouldEqual, len(tasks[0].SubTasks))
		})
	})
}
