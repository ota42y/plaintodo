package main

import (
	"testing"
	"time"
)

/*
get first task and one subtask
*/
func TestKeyValueQuery(t *testing.T) {
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

func TestNoKeyQuery(t *testing.T) {
	tasks := ReadTestTasks()

	query := NewNoKeyQuery("due", make([]Query, 0), make([]Query, 0))
	showTasks := Ls(tasks, query)

	if showTasks == nil {
		t.Errorf("filter is nil")
		t.FailNow()
	}

	if len(showTasks) != 2 {
		t.Errorf("return shuld be 2 show tasks, but %d", len(showTasks))
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("return shud be task which haven't due attribute but there is %d task", len(showTasks[0].SubTasks))
		t.FailNow()
	}

	if len(showTasks[1].SubTasks) != 0 {
		t.Errorf("return shud be task which haven't due attribute but there is %d task", len(showTasks[1].SubTasks))
		t.FailNow()
	}
}

/*
get second task and one subtask
*/
func TestBeforeDateQuery(t *testing.T) {
	tasks := ReadTestTasks()

	key := "due"
	dueTime := "2015-02-01 10:42"

	var timeformat = "2006-01-02 15:04"
	value, err := time.Parse(timeformat, dueTime)
	if err != nil {
		t.Errorf("time parse error")
		t.FailNow()
	}

	query := NewBeforeDateQuery(key, value, make([]Query, 0), make([]Query, 0))
	showTasks := Ls(tasks, query)

	if len(showTasks) == 0 {
		t.Errorf("return no tasks")
		t.FailNow()
	}

	showTask := showTasks[0]

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

	dueTime = "2015-02-02 10:42"
	value, err = time.Parse(timeformat, dueTime)
	if err != nil {
		t.Errorf("time parse error")
		t.FailNow()
	}

	query = NewBeforeDateQuery(key, value, make([]Query, 0), make([]Query, 0))
	showTasks = Ls(tasks, query)
	if len(showTasks) != 2 {
		t.Errorf("return 2 tasks but %d", len(showTasks))
		t.FailNow()
	}
}

func TestAfterDateQuery(t *testing.T) {
	tasks := ReadTestTasks()

	tasks[0].SubTasks[0].Attributes["complete"] = "2015-01-31 10:42"
	tasks[1].SubTasks[0].Attributes["complete"] = "2015-02-02 10:42"

	key := "complete"
	dueTime := "2015-02-01 00:00"

	var timeformat = "2006-01-02 15:04"
	value, err := time.Parse(timeformat, dueTime)
	if err != nil {
		t.Errorf("time parse error")
		t.FailNow()
	}

	query := NewAfterDateQuery(key, value, make([]Query, 0), make([]Query, 0))
	showTasks := Ls(tasks, query)

	if len(showTasks) == 0 {
		t.Errorf("return no tasks")
		t.FailNow()
	}

	showTask := showTasks[0]

	if showTask.Task.Name != tasks[1].Name {
		t.Errorf("filter isn't valid")
		t.FailNow()
	}

	if len(showTask.SubTasks) != 1 {
		t.Errorf("SubTasks num isn't 1")
		t.FailNow()
	}

	subTask := showTask.SubTasks[0]
	if subTask.Task.Name != tasks[1].SubTasks[0].Name {
		t.Errorf("SubTasks isn't correct")
		t.FailNow()
	}

	orQuery := make([]Query, 0)
	orQuery = append(orQuery, query)
	noCompleteOrTodayCompleteQuery := NewNoKeyQuery("due", make([]Query, 0), orQuery)
	showTasks = Ls(tasks, noCompleteOrTodayCompleteQuery)

	if len(showTasks) != 2 {
		t.Errorf("return shuld be 2 tasks but %d task", len(showTasks))
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("return shud be task which haven't due attribute but there is %d task", len(showTasks[0].SubTasks))
		t.FailNow()
	}

	// or query can reverse
	orQuery = make([]Query, 0)
	orQuery = append(orQuery, NewNoKeyQuery("due", make([]Query, 0), make([]Query, 0)))
	reverseOrQuery := NewAfterDateQuery(key, value, make([]Query, 0), orQuery)
	showTasks = Ls(tasks, reverseOrQuery)

	if len(showTasks) != 2 {
		t.Errorf("return shuld be 2 tasks but %d task", len(showTasks))
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("return shud be task which haven't due attribute but there is %d task", len(showTasks[0].SubTasks))
		t.FailNow()
	}
}

func TestSameDayQuery(t *testing.T) {
	tasks := ReadTestTasks()

	tasks[0].SubTasks[0].Attributes["complete"] = "2015-02-01 10:42"
	tasks[0].SubTasks[1].Attributes["complete"] = "2015-02-01 20:42"

	key := "complete"
	dueTime := "2015-02-01 12:00"

	var timeformat = "2006-01-02 15:04"
	value, err := time.Parse(timeformat, dueTime)
	if err != nil {
		t.Errorf("time parse error")
		t.FailNow()
	}

	query := NewSameDayQuery(key, value, make([]Query, 0), make([]Query, 0))
	showTasks := Ls(tasks, query)

	if len(showTasks) == 0 {
		t.Errorf("return no tasks")
		t.FailNow()
	}

	showTask := showTasks[0]

	if showTask.Task.Name != tasks[0].Name {
		t.Errorf("filter isn't valid")
		t.FailNow()
	}

	if len(showTask.SubTasks) != 2 {
		t.Errorf("SubTasks num isn't 2")
		t.FailNow()
	}

	subTask := showTask.SubTasks[0]
	if subTask.Task.Name != tasks[0].SubTasks[0].Name {
		t.Errorf("SubTasks isn't correct")
		t.FailNow()
	}

	if len(subTask.SubTasks) != 0 {
		t.Errorf("There is not same day query")
		t.FailNow()
	}

}
