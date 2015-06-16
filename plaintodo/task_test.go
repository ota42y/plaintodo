package main

import (
	"testing"
)

func TestReadTasks(t *testing.T) {
	tasks := ReadTestTasks()

	if len(tasks) != 2 {
		t.Errorf("there is unread task, len(tasks) should be 2 but %d", len(tasks))
		t.FailNow()
	}

	if tasks[0].Level != 0 {
		t.Errorf("top level task isn't level 0")
		t.FailNow()
	}

	if tasks[0].Id != 1 {
		t.Errorf("first task's task id shuld be 1")
		t.FailNow()
	}
}

func TestCopyTasks(t *testing.T) {
	task := &Task{
		Level:      1,
		Id:         1,
		Name:       "name",
		Attributes: make(map[string]string),
		SubTasks:   make([]*Task, 0),
	}
	task.Attributes["attr"] = "attr"

	subTask := &Task{
		Level:      2,
		Id:         2,
		Name:       "subtask",
		Attributes: make(map[string]string),
		SubTasks:   make([]*Task, 0),
	}
	subTask.Attributes["attr"] = "subattr"
	task.SubTasks = append(task.SubTasks, subTask)

	copyTask := task.Copy(3, true)

	if !task.Equal(copyTask) {
		t.Errorf("Task.Copy don't return same task")
		t.FailNow()
	}

	copyParentTask := task.Copy(5, false)
	task.SubTasks = make([]*Task, 0)
	if !task.Equal(copyParentTask) {
		t.Errorf("Task.Copy don't return same task")
		t.FailNow()
	}
}

func TestEqualTasks(t *testing.T) {
	task := &Task{
		Level:      1,
		Id:         1,
		Name:       "name",
		Attributes: make(map[string]string),
		SubTasks:   make([]*Task, 0),
	}
	task.Attributes["attr"] = "attr"

	subTask := &Task{
		Level:      2,
		Id:         2,
		Name:       "subtask",
		Attributes: make(map[string]string),
		SubTasks:   make([]*Task, 0),
	}
	subTask.Attributes["attr"] = "subattr"
	task.SubTasks = append(task.SubTasks, subTask)

	eqTask := &Task{
		Level:      1,
		Id:         3,
		Name:       "name",
		Attributes: make(map[string]string),
		SubTasks:   make([]*Task, 0),
	}
	task.Attributes["attr"] = "attr"

	eqSubTask := &Task{
		Level:      2,
		Id:         4,
		Name:       "subtask",
		Attributes: make(map[string]string),
		SubTasks:   make([]*Task, 0),
	}
	eqSubTask.Attributes["attr"] = "subattr"
	eqTask.SubTasks = append(eqTask.SubTasks, eqSubTask)

	eqTask.Name = "notEq"
	if task.Equal(eqTask) {
		t.Errorf("even task name isn't equal, Task.Equal return true")
		t.FailNow()
	}
	eqTask.Name = "name"

	eqTask.Level = 10
	if task.Equal(eqTask) {
		t.Errorf("even task level isn't equal, Task.Equal return true")
		t.FailNow()
	}
	eqTask.Level = 1

	eqTask.Attributes["attr"] = "aaaa"
	if task.Equal(eqTask) {
		t.Errorf("even task attribute isn't equal, Task.Equal return true")
		t.FailNow()
	}
	eqTask.Attributes["attr"] = "attr"

	eqTask.Attributes["bbb"] = "test"
	if task.Equal(eqTask) {
		t.Errorf("even task attribute isn't equal, Task.Equal return true")
		t.FailNow()
	}
	delete(eqTask.Attributes, "bbb")

	eqTask.SubTasks = make([]*Task, 0)
	if task.Equal(eqTask) {
		t.Errorf("even subtask num isn't equal, Task.Equal return true")
		t.FailNow()
	}
	eqTask.SubTasks = append(eqTask.SubTasks, eqSubTask)

	eqSubTask.Name = "notEq"
	if task.Equal(eqTask) {
		t.Errorf("even subtask isn't equal, Task.Equal return true")
		t.FailNow()
	}
	eqSubTask.Name = "subtask"

	if !task.Equal(eqTask) {
		t.Errorf("task isn't equal")
		t.FailNow()
	}
}

func TestCreateSubTasks(t *testing.T) {
	tasks := ReadTestTasks()

	if len(tasks) != 2 {
		t.Errorf("read top level subtasks failed, len(tasks) shuld be 2 but %d", len(tasks))
		t.FailNow()
	}

	task := tasks[0]

	if len(task.SubTasks) != 2 {
		t.Errorf("read subtasks failed, SubTasks count shuld be 2 but %d", len(task.SubTasks))
		t.FailNow()
	}

	subTask := task.SubTasks[1]

	if subTask.Level != 1 {
		t.Errorf("read subtask's data failed %v", subTask)
		t.FailNow()
	}

	if subTask.Name != "buy items" {
		t.Errorf("read subtask's data failed %v", subTask)
		t.FailNow()
	}

	if subTask.Id != 4 {
		t.Errorf("%s's id shuld be 4 but %d", subTask.Name, subTask.Id)
		t.FailNow()
	}

	if len(subTask.SubTasks) != 3 {
		t.Errorf("read subtask's subtask failed")
		t.FailNow()
	}

	subSubTask := subTask.SubTasks[0]
	if subSubTask.Level != 2 {
		t.Errorf("read subtask's subtask level failed")
		t.FailNow()
	}

	if subSubTask.Id != 5 {
		t.Errorf("%s's id shuld be 5 but %d", subSubTask.Name, subSubTask.Id)
		t.FailNow()
	}
}

func TestNewTask(t *testing.T) {
	line := "    add music to player"
	task, err := NewTask(line, 1)

	if err != nil {
		t.Errorf("NewTask return error %v", err)
		t.FailNow()
	}

	if task == nil {
		t.Errorf("task is nil")
		t.FailNow()
	}

	if task.Level != 2 {
		t.Errorf("task.Level shuold be 3 but %d", task.Level)
		t.FailNow()
	}

	if task.Id != 1 {
		t.Errorf("task.Id shuold be 1 but %d", task.Id)
		t.FailNow()
	}

	if len(task.Attributes) != 0 {
		t.Errorf("task.Attributes shuld be empty but %v", task.Attributes)
		t.FailNow()
	}

	if task.Attributes == nil {
		t.Errorf("task.Attributes shuldn't be nil")
		t.FailNow()
	}

	correctName := "add music to player"
	if task.Name != correctName {
		t.Errorf("task.Name shuold be %s but %s", correctName, task.Name)
		t.FailNow()
	}

	taskString := task.String(false)
	if taskString != line {
		t.Errorf("task.String return invalid string %s", taskString)
		t.FailNow()
	}

	taskString = task.String(true)
	if taskString != "    add music to player :id 1" {
		t.Errorf("task.String return invalid string %s", taskString)
		t.FailNow()
	}
}

func TestNewTaskWithAttributes(t *testing.T) {
	line := "    create a set list :important :repeat every 1 day :start 2015-02-01 :url http://ota42y.com"
	task, err := NewTask(line, 1)

	if err != nil {
		t.Errorf("NewTask return error %v", err)
		t.FailNow()
	}

	if task == nil {
		t.Errorf("task is nil")
		t.FailNow()
	}

	if task.Level != 2 {
		t.Errorf("task.Level shuold be 2 but %d", task.Level)
		t.FailNow()
	}

	if task.Id != 1 {
		t.Errorf("task.Id shuold be 1 but %d", task.Id)
		t.FailNow()
	}

	correctName := "create a set list"
	if task.Name != correctName {
		t.Errorf("task.Name shuold be '%s' but '%s'", correctName, task.Name)
		t.FailNow()
	}

	// :url http://ota42y.com :start 2015-02-01 :repeat every 1 day"
	attributes := make(map[string]string)
	attributes["url"] = "http://ota42y.com"
	attributes["start"] = "2015-02-01"
	attributes["repeat"] = "every 1 day"
	attributes["important"] = ""

	for key, value := range attributes {
		if task.Attributes[key] != value {
			t.Errorf("key: %s shuld be %s but %s", key, value, task.Attributes[key])
			t.FailNow()
		}
	}

	if len(attributes) != len(task.Attributes) {
		t.Errorf("Task.Attributes shuld be %d num, but %d", len(attributes), len(task.Attributes))
		t.FailNow()
	}

	taskString := task.String(false)
	if taskString != line {
		t.Errorf("task.String return invalid string %s", taskString)
		t.FailNow()
	}

	taskString = task.String(true)
	correctString := "    create a set list :id 1 :important :repeat every 1 day :start 2015-02-01 :url http://ota42y.com"
	if taskString != correctString {
		t.Errorf("task.String shuld return '%s' string '%s'", correctString, taskString)
		t.FailNow()
	}
}

func TestNewTaskError(t *testing.T) {
	line := "    "
	task, err := NewTask(line, 1)

	if err == nil {
		t.Errorf("blank line return err, but err is nil")
		t.FailNow()
	}

	if task != nil {
		t.Errorf("when error return, task shuld be nil, but %v", task)
		t.FailNow()
	}

	correctName := "blank line"
	if err.Error() != correctName {
		t.Errorf("task.Name shuold be '%s' but '%s'", correctName, task.Name)
		t.FailNow()
	}
}

func TestGetTask(t *testing.T) {
	tasks := ReadTestTasks()

	task := GetTask(6, tasks)
	if task == nil {
		t.Errorf("GetTask shuld return Task.Id = 6 task, but nil")
		t.FailNow()
	}

	if task.Id != 6 {
		t.Errorf("GetTask shuld return Task.Id = 6 task, but other task return")
		t.FailNow()
	}

	task = GetTask(0, tasks)
	if task != nil {
		t.Errorf("GetTask shuld return nil when task isn't exist, but %v", task)
		t.FailNow()
	}
}
