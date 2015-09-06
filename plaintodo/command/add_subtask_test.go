package command

import (
	"fmt"
	"testing"

	"../util"
)

func TestAddSubTaskCommand(t *testing.T) {
	taskName := "create sub task"
	taskStart := "2015-02-01"

	addTask := NewAddTask()
	addSubTask := NewAddSubTask()

	config, buf := util.ReadTestConfigRelativePath("..")
	s := &State{
		Config: config,
	}
	s.Tasks = util.ReadTestTaskRelativePath("../")

	taskID := s.MaxTaskID

	input := taskName + " :id 6 :start " + taskStart
	addSubTask.Execute(input, s)

	parent := s.Tasks[0].SubTasks[1].SubTasks[1]

	if len(parent.SubTasks) == 0 {
		t.Errorf("SubTask not add\n%s", buf.String())
		t.FailNow()
	}

	task := parent.SubTasks[0]
	if task.Level != parent.Level+1 {
		t.Errorf("Subtask level shuld be %d but %d", parent.Level+1, task.Level)
		t.FailNow()
	}

	if task.Name != taskName {
		t.Errorf("Task name shud %s, but %s", taskName, task.Name)
		t.FailNow()
	}

	if task.Attributes["start"] != taskStart {
		t.Errorf("Task start shud %s, but %s", taskStart, task.Attributes["start"])
		t.FailNow()
	}

	if taskID+1 != task.ID {
		t.Errorf("Task's id shud be %d but %d", taskID+1, task.ID)
		t.FailNow()
	}

	if s.MaxTaskID != task.ID {
		t.Errorf("Automaton.MaxTaskID shuld be %d, but %d", s.MaxTaskID, task.ID)
		t.FailNow()
	}
	taskID = s.MaxTaskID

	outputString := buf.String()
	correctString := "Create SubTask:\nParent: " + parent.String(true) + "\nSubTask: " + task.String(true) + "\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	buf.Reset()
	addSubTask.Execute("test", s)

	outputString = buf.String()
	correctString = "There're no parent id"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	buf.Reset()
	addSubTask.Execute("parent task not exit :id 0", s)

	outputString = buf.String()
	correctString = "Create SubTask error: thee is no task which have :id 0\n"
	if outputString != correctString {
		t.Errorf("Output %s, but %s", correctString, outputString)
		t.FailNow()
	}

	if s.MaxTaskID != taskID {
		t.Errorf("When error occerd, Automaton.MaxTaskID shuldn't change but %d", taskID)
		t.FailNow()
	}

	addTask.Execute("task test11", s)
	addTask.Execute("task test12", s)

	parentID := s.MaxTaskID
	input = fmt.Sprintf("child in task12 :id %d", parentID)
	addSubTask.Execute(input, s)

	lastPos := len(s.Tasks) - 1
	if len(s.Tasks[lastPos].SubTasks) == 1 {
		t.Errorf("Subtask isn't added by id=%d task's child", parentID)
		t.FailNow()
	}
}
