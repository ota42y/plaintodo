package command

import (
	"testing"
	"time"

	"../util"
)

func TestStart(t *testing.T) {
	cmd := NewStart()

	config, buf := util.ReadTestConfigRelativePath("../")

	tasks := util.ReadTestTaskRelativePath("../")
	s := &State{
		Tasks:  tasks,
		Config: config,
	}

	now := time.Now()

	task := s.Tasks[0].SubTasks[1].SubTasks[0]
	if _, ok := task.Attributes["start"]; ok {
		t.Errorf("task already set start attribute, test data is invalid %v", task)
		t.FailNow()
	}

	terminate := cmd.Execute(" :id 5", s)
	if terminate {
		t.Errorf("Start.Execute shud be return false")
		t.FailNow()
	}

	value, ok := task.Attributes["start"]
	if !ok {
		t.Errorf("start attribute not set")
		t.FailNow()
	}

	dateTime, ok := util.ParseTime(value)
	diff := dateTime.Sub(now)
	if diff.Minutes() < -2 || 2 < diff.Minutes() {
		t.Errorf("set time (%v) isn't now because %v minutes after", value, diff.Seconds())
		t.FailNow()
	}

	task.Attributes["start"] = time.Now().AddDate(1, 0, 0).Format(util.DateTimeFormat)

	terminate = cmd.Execute(" :id 5", s)
	dateTime, ok = util.ParseTime(value)
	diff = dateTime.Sub(now)
	if diff.Minutes() < -2 || 2 < diff.Minutes() {
		t.Errorf("set new start time, but old time isn't overwrited")
		t.FailNow()
	}

	s.Tasks = util.ReadTestTaskRelativePath("../")
	task = s.Tasks[0].SubTasks[1].SubTasks[0]
	task.Attributes["lock"] = ""

	buf.Reset()
	cmd.Execute(" :id 5", s)
	_, ok = task.Attributes["start"]
	if ok {
		t.Errorf("lock task mustn't change, but attribute set")
		t.FailNow()
	}

	if buf.String() != "Task :id 5 is locked\n" {
		t.Errorf("error message in invalid '%s'", buf.String())
		t.FailNow()
	}
}
