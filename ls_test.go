package main

import "testing"

func TestLs(t *testing.T) {
	filename := "task.txt"
	tasks := ReadTasks(filename)

	showTasks := Ls(tasks)
	if len(showTasks) != 2 {
		t.Errorf("top level shud be 2 but %d", len(showTasks))
		t.FailNow()
	}

	subTasks := showTasks[0].SubTasks
	if len(subTasks) != 2 {
		t.Errorf("top level shud be 2 but %d", len(subTasks))
		t.FailNow()
	}

	subTasks = subTasks[1].SubTasks
	if len(subTasks) != 3 {
		t.Errorf("top level shud be 3 but %d", len(subTasks))
		t.FailNow()
	}
}

/*
get first task and all sub tasks
 */
func TestFilter(t *testing.T){
	filename := "task.txt"
	tasks := ReadTasks(filename)

	showTask := filter(tasks[0])
	if showTask == nil{
		t.Errorf("filter is nil")
		t.FailNow()
	}

	if showTask.Task.Name != tasks[0].Name {
		t.Errorf("filter isn't valid")
		t.FailNow()
	}

	if len(showTask.SubTasks) != 2{
		t.Errorf("SubTasks num isn't 2")
		t.FailNow()
	}
}
