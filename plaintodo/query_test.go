package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"

	"./ls"
	"./query"
	"./util"
)

/*
get first task and one subtask
*/
func TestKeyValueQuery(t *testing.T) {
	tasks := util.ReadTestTasks()

	Convey("correct", t, func() {
		q := query.NewKeyValue("start", "2015-01-31", query.CreateBlankQueryArray(), query.CreateBlankQueryArray())
		showTask := ls.Filter(tasks[0], q)

		So(showTask, ShouldNotBeNil)

		Convey("top level task", func() {
			topLevelTask := showTask.Task
			So(topLevelTask.Name, ShouldEqual, tasks[0].Name)
		})

		Convey("sub task", func() {
			So(len(showTask.SubTasks), ShouldEqual, 1)

			subTask := showTask.SubTasks[0]
			So(subTask.Task.Name, ShouldEqual, tasks[0].SubTasks[0].Name)
		})

		Convey("sub task name valid", nil)

	})

	Convey("error", t, func() {
		Convey("no hit", func() {
			q := query.NewKeyValue("start", "2015-01-31", query.CreateBlankQueryArray(), query.CreateBlankQueryArray())
			showTask := ls.Filter(tasks[1], q)
			So(showTask, ShouldBeNil)
		})
	})
}

func TestNoKeyQuery(t *testing.T) {
	tasks := util.ReadTestTasks()

	q := query.NewNoKey("start", make([]query.Query, 0), make([]query.Query, 0))
	showTasks := ls.Ls(tasks, q)

	if showTasks == nil {
		t.Errorf("filter is nil")
		t.FailNow()
	}

	if len(showTasks) != 2 {
		t.Errorf("return shuld be 2 show tasks, but %d", len(showTasks))
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("return shud be task which haven't start attribute but there is %d task", len(showTasks[0].SubTasks))
		t.FailNow()
	}

	if len(showTasks[1].SubTasks) != 0 {
		t.Errorf("return shud be task which haven't start attribute but there is %d task", len(showTasks[1].SubTasks))
		t.FailNow()
	}
}

/*
get second task and one subtask
*/
func TestBeforeDateQuery(t *testing.T) {
	tasks := util.ReadTestTasks()

	key := "start"
	startTime := "2015-02-01 10:42"

	var timeformat = "2006-01-02 15:04"
	value, err := time.Parse(timeformat, startTime)
	if err != nil {
		t.Errorf("time parse error")
		t.FailNow()
	}

	q := query.NewBeforeDate(key, value, make([]query.Query, 0), make([]query.Query, 0))
	showTasks := ls.Ls(tasks, q)

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

	startTime = "2015-02-02 10:42"
	value, err = time.Parse(timeformat, startTime)
	if err != nil {
		t.Errorf("time parse error")
		t.FailNow()
	}

	q = query.NewBeforeDate(key, value, make([]query.Query, 0), make([]query.Query, 0))
	showTasks = ls.Ls(tasks, q)
	if len(showTasks) != 2 {
		t.Errorf("return 2 tasks but %d", len(showTasks))
		t.FailNow()
	}
}

func TestAfterDateQuery(t *testing.T) {
	tasks := util.ReadTestTasks()

	tasks[0].SubTasks[0].Attributes["complete"] = "2015-01-31 10:42"
	tasks[1].SubTasks[0].Attributes["complete"] = "2015-02-02 10:42"

	key := "complete"
	startTime := "2015-02-01 00:00"

	var timeformat = "2006-01-02 15:04"
	value, err := time.Parse(timeformat, startTime)
	if err != nil {
		t.Errorf("time parse error")
		t.FailNow()
	}

	q := NewAfterDateQuery(key, value, make([]query.Query, 0), make([]query.Query, 0))
	showTasks := ls.Ls(tasks, q)

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

	orQuery := make([]query.Query, 0)
	orQuery = append(orQuery, q)
	noCompleteOrTodayCompleteQuery := query.NewNoKey("start", make([]query.Query, 0), orQuery)
	showTasks = ls.Ls(tasks, noCompleteOrTodayCompleteQuery)

	if len(showTasks) != 2 {
		t.Errorf("return shuld be 2 tasks but %d task", len(showTasks))
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("return shud be task which haven't start attribute but there is %d task", len(showTasks[0].SubTasks))
		t.FailNow()
	}

	// or query can reverse
	orQuery = make([]query.Query, 0)
	orQuery = append(orQuery, query.NewNoKey("start", make([]query.Query, 0), make([]query.Query, 0)))
	reverseOrQuery := NewAfterDateQuery(key, value, make([]query.Query, 0), orQuery)
	showTasks = ls.Ls(tasks, reverseOrQuery)

	if len(showTasks) != 2 {
		t.Errorf("return shuld be 2 tasks but %d task", len(showTasks))
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("return shud be task which haven't start attribute but there is %d task", len(showTasks[0].SubTasks))
		t.FailNow()
	}
}

func TestSameDayQuery(t *testing.T) {
	tasks := util.ReadTestTasks()

	tasks[0].SubTasks[0].Attributes["complete"] = "2015-02-01 10:42"
	tasks[0].SubTasks[1].Attributes["complete"] = "2015-02-01 20:42"

	key := "complete"
	startTime := "2015-02-01 12:00"

	var timeformat = "2006-01-02 15:04"
	value, err := time.Parse(timeformat, startTime)
	if err != nil {
		t.Errorf("time parse error")
		t.FailNow()
	}

	q := NewSameDayQuery(key, value, make([]query.Query, 0), make([]query.Query, 0))
	showTasks := ls.Ls(tasks, q)

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
