package main

import (
	"testing"
)

func TestReadTasks(t *testing.T) {
	filename := "task.txt"
	tasks := ReadTasks(filename)

	if len(tasks) != 2{
		t.Errorf("there is unread task, len(tasks) should be 2 but %d", len(tasks))
		t.FailNow()
	}

	if tasks[0].Level != 0{
		t.Errorf("top level task isn't level 0")
		t.FailNow()
	}
}

func TestCreateSubTasks(t *testing.T) {
	filename := "task.txt"
	tasks := ReadTasks(filename)

	if len(tasks) != 2{
		t.Errorf("read top level subtasks failed, len(tasks) shuld be 2 but %d", len(tasks))
		t.FailNow()
	}

	task := tasks[0]

	if len(task.SubTasks) != 2{
		t.Errorf("read subtasks failed, SubTasks count shuld be 2 but %d", len(task.SubTasks))
		t.FailNow()
	}

	subTask := task.SubTasks[1]

	if subTask.Level != 1{
		t.Errorf("read subtask's data failed %v", subTask)
		t.FailNow()
	}

	if subTask.Name != "buy items"{
		t.Errorf("read subtask's data failed %v", subTask)
		t.FailNow()
	}

	if len(subTask.SubTasks) != 3{
		t.Errorf("read subtask's subtask failed")
		t.FailNow()
	}

	subSubTask := subTask.SubTasks[0]
	if subSubTask.Level != 2{
		t.Errorf("read subtask's subtask level failed")
		t.FailNow()
	}
}

func TestNewTask(t *testing.T) {
	line := "    add music to player"
	task := NewTask(line)

	if task == nil{
		t.Errorf("task is nil")
		t.FailNow()
	}

	if task.Level != 2{
		t.Errorf("task.Level shuold be 3 but %d", task.Level)
		t.FailNow()
	}

	correctName := "add music to player"
	if task.Name != correctName{
		t.Errorf("task.Name shuold be %s but %s", correctName, task.Name)
		t.FailNow()
	}
}

func TestNewTaskWithAttributes(t *testing.T) {
	line := "    create a set list :due 2015-01-31 :repeat daily"
	task := NewTask(line)

	if task == nil{
		t.Errorf("task is nil")
		t.FailNow()
	}

	if task.Level != 2{
		t.Errorf("task.Level shuold be 2 but %d", task.Level)
		t.FailNow()
	}

	correctName := "create a set list"
	if task.Name != correctName{
		t.Errorf("task.Name shuold be '%s' but '%s'", correctName, task.Name)
		t.FailNow()
	}
}
