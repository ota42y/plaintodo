package main

import (
	"testing"
)

/*
get first task and one subtask
*/
func TestFilterWithDueDate(t *testing.T) {
	tasks := ReadTestTasks()

	query := NewKeyValueQuery("due", "2015-01-31", make([]Query, 0), make([]Query, 0))
	showTask := filter(tasks[0], query)
	if showTask == nil {
		t.Errorf("filter is nil")
		t.FailNow()
	}

	if showTask.Task.Name != tasks[0].Name {
		t.Errorf("filter isn't valid")
		t.FailNow()
	}

	if len(showTask.SubTasks) != 1 {
		t.Errorf("SubTasks num isn't 1")
		t.FailNow()
	}

	subTask := showTask.SubTasks[0]
	if subTask.Task.Name != tasks[0].SubTasks[0].Name {
		t.Errorf("SubTasks isn't correct")
		t.FailNow()
	}

	showTask = filter(tasks[1], query)
	if showTask != nil {
		t.Errorf("if not match any query, get nil but get %v", showTask)
		t.FailNow()
	}
}