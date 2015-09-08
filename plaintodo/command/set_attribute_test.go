package command

import (
	"testing"

	"../util"
)

func TestSetAttributeCommand(t *testing.T) {
	cmd := NewSetAttribute()
	url := "http://example.com"

	config, buf := util.ReadTestConfigRelativePath("../")

	tasks := util.ReadTestTaskRelativePath("../")

	s := &State{
		Tasks:  tasks,
		Config: config,
	}

	terminate := cmd.Execute(":url "+url, s)
	if terminate {
		t.Errorf("SetAttributeCommand.Execute shud be return false")
		t.FailNow()
	}

	outputString := buf.String()
	correctString := "not exist :id\n"
	if outputString != correctString {
		t.Errorf("Shuld output '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}

	task := s.Tasks[0].SubTasks[0].SubTasks[0]
	task.Attributes["url"] = url
	correctString = "set attribute " + task.String(true) + "\n"
	delete(task.Attributes, "url")

	buf.Reset()
	terminate = cmd.Execute(":id 3 :url "+url, s)
	if terminate {
		t.Errorf("SetAttributeCommand.Execute shud be return false")
		t.FailNow()
	}

	value, ok := task.Attributes["url"]
	if !ok {
		t.Errorf("attribute not set")
		t.FailNow()
	}

	if value != url {
		t.Errorf("set attribute shuld %s, but %s", url, value)
		t.FailNow()
	}

	outputString = buf.String()
	if outputString != correctString {
		t.Errorf("Shuld output \n'%s', but \n'%s'", correctString, outputString)
		t.FailNow()
	}

	buf.Reset()
	terminate = cmd.Execute(":id 0 :url "+url, s)

	outputString = buf.String()
	correctString = "there is no exist :id 0 task\n"
	if outputString != correctString {
		t.Errorf("Shuld output '%s', but '%s'", correctString, outputString)
		t.FailNow()
	}

	buf.Reset()
	s.Tasks = util.ReadTestTaskRelativePath("../")

	task = s.Tasks[0].SubTasks[0].SubTasks[0]
	task.Attributes["lock"] = ""
	terminate = cmd.Execute(":id 3 :url "+url, s)

	_, ok = task.Attributes["url"]
	if ok {
		t.Errorf("lock task mustn't change, but attribute set")
		t.FailNow()
	}

	if buf.String() != "Task :id 3 is locked\n" {
		t.Errorf("error message in invalid '%s'", buf.String())
		t.FailNow()
	}
}
