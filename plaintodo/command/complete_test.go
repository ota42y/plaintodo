package command

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"

	"../task"
	"../util"
	"fmt"
)

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

func TestCompleteTaskExecute(t *testing.T) {
	cmd := NewComplete()
	config, buf := util.ReadTestConfigRelativePath("..")
	s := &State{
		Config: config,
	}

	Convey("correct", t, func() {
		Convey("complete", func() {
			s.Tasks = util.ReadTestTaskRelativePath("../")
			option := ":id 1"

			buf.Reset()
			cmd.Execute(option, s)
			_, ok := s.Tasks[0].Attributes["complete"]
			So(ok, ShouldBeTrue)

			correctString := fmt.Sprintf("Complete %s and %d sub tasks\n", s.Tasks[0].Name, 6)
			So(correctString, ShouldEqual, buf.String())
		})
	})

	Convey("incorrect", t, func() {
		Convey("not set", func() {
			s.Tasks = util.ReadTestTaskRelativePath("../")
			option := ":id 42"

			buf.Reset()
			cmd.Execute(option, s)
			_, no := s.Tasks[0].Attributes["complete"]
			So(no, ShouldBeFalse)

			correctString := fmt.Sprintf("There is no Task which have task id: %d\n", 42)
			So(correctString, ShouldEqual, buf.String())
		})

		Convey("lock task", func() {
			s.Tasks = util.ReadTestTaskRelativePath("../")
			s.Tasks[0].Attributes["lock"] = ""

			option := ":id 1"

			buf.Reset()
			cmd.Execute(option, s)
			_, no := s.Tasks[0].Attributes["complete"]
			So(no, ShouldBeFalse)

			correctString := fmt.Sprintf("task :id 1 is locked")
			So(correctString, ShouldEqual, buf.String())
		})
	})
}

func TestCompleteTask(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("../")
	cmd := NewComplete()

	result, err := cmd.completeTask(0, tasks)
	if result != nil {
		t.Errorf("If there is no task with taskID, completeTask shuld return nil, but %v", result.completeTask)
		t.FailNow()
	}

	alreadyCompleted := "2014-01-01"
	tasks[0].SubTasks[1].SubTasks[1].Attributes["complete"] = alreadyCompleted

	result, err = cmd.completeTask(4, tasks)
	if result.completeTask == nil {
		t.Errorf("If there is task with taskID, completeTask shuld return complete task, but nil")
		t.FailNow()
	}

	if len(tasks) != 2 {
		t.Errorf("task num shudn't change")
		t.FailNow()
	}

	if result.completeNum != 3 {
		t.Errorf("If there is task with taskID, completeTask shuld return complete subtask num (4) but %d", result.completeNum)
		t.FailNow()
	}

	if !isAllCompleted(tasks[0].SubTasks[1]) {
		t.Errorf("not complete selected task and all sub tasks")
		t.FailNow()
	}

	completeString := tasks[0].SubTasks[1].Attributes["complete"]
	_, err = time.Parse(util.DateTimeFormat, completeString)
	if err != nil {
		t.Errorf("Task complete format invalid '%s'", completeString)
		t.FailNow()
	}

	completeString = tasks[0].SubTasks[1].SubTasks[0].Attributes["complete"]
	_, err = time.Parse(util.DateTimeFormat, completeString)
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

func TestSetNewRepeat(t *testing.T) {
	cmd := NewComplete()
	tk := &task.Task{
		Attributes: make(map[string]string),
	}

	now := time.Now()
	base := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.Local)
	tk.Attributes["start"] = base.Format(util.DateTimeFormat)

	tk.Attributes["repeat"] = "every 1 day"
	cmd.setNewRepeat(now, tk)
	correct := base.AddDate(0, 0, 1)
	correctString := correct.Format(util.DateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(util.DateTimeFormat)
	tk.Attributes["repeat"] = "every 1 month"
	cmd.setNewRepeat(now, tk)
	correct = base.AddDate(0, 1, 0)
	correctString = correct.Format(util.DateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(util.DateTimeFormat)
	tk.Attributes["repeat"] = "every 1 year"
	cmd.setNewRepeat(now, tk)
	correct = base.AddDate(1, 0, 0)
	correctString = correct.Format(util.DateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(util.DateTimeFormat)
	tk.Attributes["repeat"] = "every 2 week"
	cmd.setNewRepeat(now, tk)
	correct = base.AddDate(0, 0, 14)
	correctString = correct.Format(util.DateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(util.DateTimeFormat)
	tk.Attributes["repeat"] = "every 30 minutes"
	cmd.setNewRepeat(now, tk)
	correct = base.Add(30 * time.Minute)
	correctString = correct.Format(util.DateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(util.DateTimeFormat)
	tk.Attributes["repeat"] = "every 2 hour"
	cmd.setNewRepeat(now, tk)
	correct = base.Add(2 * time.Hour)
	correctString = correct.Format(util.DateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}

	tk.Attributes["start"] = base.Format(util.DateTimeFormat)
	tk.Attributes["repeat"] = "after 4 day"
	cmd.setNewRepeat(now, tk)
	correct = now.AddDate(0, 0, 4)
	correctString = correct.Format(util.DateTimeFormat)

	if correctString != tk.Attributes["start"] {
		t.Errorf("Time shuld be %v but %v", correctString, tk.Attributes["start"])
		t.FailNow()
	}
}

func TestCompleteRepeatTask(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("../")
	baseTask := tasks[1].SubTasks[0]

	postponeCommand := NewPostpone()
	optionMap := make(map[string]string)
	optionMap["postpone"] = "1 day"
	postponeCommand.Postpone(baseTask, optionMap)

	cmd := NewComplete()

	result, _ := cmd.completeTask(8, tasks)
	if result.completeTask == nil {
		t.Errorf("If there is task with taskID, completeTask shuld return complete task, but nil")
		t.FailNow()
	}

	if result.completeNum != 2 {
		t.Errorf("If there is task with taskID, completeTask shuld return complete subtask num (2) but %d", result.completeNum)
		t.FailNow()
	}

	newTasks := result.newTasks
	if len(newTasks) != 3 {
		t.Errorf("If repeat task complete, task will copy")
		t.FailNow()
	}

	baseStart, baseOk := util.ParseTime(newTasks[1].SubTasks[0].Attributes["start"])
	repeatStart, repeatOk := util.ParseTime(newTasks[2].SubTasks[0].Attributes["start"])
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
