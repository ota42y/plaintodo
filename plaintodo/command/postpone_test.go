package command

import (
	"fmt"
	"testing"
	"time"

	"../util"
)

func TestPostponeCommand(t *testing.T) {
	cmd := NewPostpone()

	config, buf := util.ReadTestConfigRelativePath("../")
	tasks := util.ReadTestTaskRelativePath("../")

	s := &State{
		Tasks:  tasks,
		Config: config,
	}

	now := time.Now()

	tk := s.Tasks[0].SubTasks[1].SubTasks[0]
	if _, ok := tk.Attributes["postpone"]; ok {
		t.Errorf("task already set postpone attribute, test data is invalid %v", tk)
		t.FailNow()
	}

	cmd.Execute(":id 5 :postpone 1 month", s)
	outputString := buf.String()
	correctString := fmt.Sprintln("task :id", tk.ID, "haven't start attribute, so postpone not work")
	if outputString != correctString {
		t.Errorf("shuld return '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	tk.Attributes["start"] = "test"
	cmd.Execute(":id 5 :postpone 1 month", s)
	outputString = buf.String()
	correctString = fmt.Sprintln("test is invalid format, so postpone not work")
	if outputString != correctString {
		t.Errorf("shuld return '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}
	buf.Reset()

	// set start 0 time
	tk.Attributes["start"] = time.Unix(0, 0).Format(util.DateTimeFormat)

	// invalid case
	cmd.Execute(":id 5 :postpone 1", s)
	outputString = buf.String()
	correctString = fmt.Sprintln("1 is invalid format")
	if outputString != correctString {
		t.Errorf("shuld return '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}

	terminate := cmd.Execute(":id 5 :postpone 1 month", s)
	if terminate {
		t.Errorf("PostPoneCommand.Execute shud be return false")
		t.FailNow()
	}

	value, ok := tk.Attributes["postpone"]
	if !ok {
		t.Errorf("postpone attribute not set %v", tk)
		t.FailNow()
	}

	dateTime, ok := util.ParseTime(value)
	if !ok {
		t.Errorf("postpone attribute value is invalid formt %s", value)
		t.FailNow()
	}

	diff := dateTime.Sub(now.AddDate(0, 1, 0))
	if diff.Minutes() < -2 || 2 < diff.Minutes() {
		t.Errorf("postpone time (%v) isn't 1 month ofter because %v minutes after", value, diff.Minutes())
		t.FailNow()
	}

	buf.Reset()
	s.Tasks = util.ReadTestTaskRelativePath("../")
	tk = s.Tasks[0].SubTasks[1].SubTasks[0]

	tk.Attributes["start"] = time.Unix(0, 0).Format(util.DateTimeFormat)
	tk.Attributes["lock"] = ""

	cmd.Execute(":id 5 :postpone 1 month", s)

	_, ok = tk.Attributes["postpone"]
	if ok {
		t.Errorf("lock task mustn't change, but postpone attribute set")
		t.FailNow()
	}

	if buf.String() != fmt.Sprintf("Task :id 5 is locked\n") {
		t.Errorf("error message in invalid '%s'", buf.String())
		t.FailNow()
	}
}
