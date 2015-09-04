package command

import (
	"bytes"
	"testing"

	"../util"
)

func TestAddTaskCommand(t *testing.T) {
	taskName := "create new task"
	taskStart := "2015-02-01"

	addTaskCommand := NewAddTask()

	config, buf := util.ReadTestConfigRelativePath("..")
	s := &State{
		Config: config,
	}

	input := taskName + " :start " + taskStart
	addTaskCommand.Execute(input, s)

	if len(s.Tasks) == 0 {
		t.Errorf("Task not add")
		t.FailNow()
	}

	task := s.Tasks[0]
	if task.Name != taskName {
		t.Errorf("Task name shud %s, but %s", taskName, task.Name)
		t.FailNow()
	}

	if task.Attributes["start"] != taskStart {
		t.Errorf("Task start shud %s, but %s", taskStart, task.Attributes["start"])
		t.FailNow()
	}

	if s.MaxTaskID != task.ID {
		t.Errorf("Automaton.MaxTaskID shuld be %d, but %d", s.MaxTaskID, task.ID)
		t.FailNow()
	}

	outputString := buf.String()
	correctString := "Create task: " + taskName + " :id 1 :start " + taskStart + "\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	taskID := s.MaxTaskID

	buf = &bytes.Buffer{}
	config.Writer = buf
	addTaskCommand.Execute("", s)

	outputString = buf.String()
	correctString = "Create task error: blank line\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	if s.MaxTaskID != taskID {
		t.Errorf("When error occerd, Automaton.MaxTaskID shuldn't change but %d", taskID)
		t.FailNow()
	}
}
