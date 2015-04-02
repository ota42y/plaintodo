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

	if len(subTask.SubTasks) != 3 {
		t.Errorf("read subtask's subtask failed")
		t.FailNow()
	}

	subSubTask := subTask.SubTasks[0]
	if subSubTask.Level != 2 {
		t.Errorf("read subtask's subtask level failed")
		t.FailNow()
	}
}

func TestNewTask(t *testing.T) {
	line := "    add music to player"
	task, err := NewTask(line)

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

	correctName := "add music to player"
	if task.Name != correctName {
		t.Errorf("task.Name shuold be %s but %s", correctName, task.Name)
		t.FailNow()
	}

	taskString := task.String()
	if taskString != line {
		t.Errorf("task.String return invalid string %s", taskString)
		t.FailNow()
	}
}

func TestNewTaskWithAttributes(t *testing.T) {
	line := "    create a set list :due 2015-02-01 :important :repeat every 1 day :url http://ota42y.com"
	task, err := NewTask(line)

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

	correctName := "create a set list"
	if task.Name != correctName {
		t.Errorf("task.Name shuold be '%s' but '%s'", correctName, task.Name)
		t.FailNow()
	}

	// :url http://ota42y.com :due 2015-02-01 :repeat every 1 day"
	attributes := make(map[string]string)
	attributes["url"] = "http://ota42y.com"
	attributes["due"] = "2015-02-01"
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

	taskString := task.String()
	if taskString != line {
		t.Errorf("task.String return invalid string %s", taskString)
		t.FailNow()
	}
}

func TentNewTaskError(t *testing.T) {
	line := "    "
	task, err := NewTask(line)

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
