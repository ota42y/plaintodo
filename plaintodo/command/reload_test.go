package command

import (
	"bufio"
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"../task"
	"../util"
)

func TestReloadCommand(t *testing.T) {
	cmd := NewReload()

	config, _ := util.ReadTestConfigRelativePath("../")
	config.Task.DefaultFilename = "test_task.txt"
	config.Task.TaskFolder = "../"

	s := &State{
		Config: config,
	}

	terminate := cmd.Execute("", s)
	if terminate {
		t.Errorf("ReloadCommand.Execute shud be return false")
		t.FailNow()
	}

	if len(s.Tasks) == 0 {
		t.Errorf("Task num shuldn't be 0")
		t.FailNow()
	}

	id := s.Tasks[1].SubTasks[0].ID
	if s.MaxTaskID != id {
		t.Errorf("Save max task id %d, but %d", id, s.MaxTaskID)
		t.FailNow()
	}
}

func createSubTaskTestState() (*State, *bytes.Buffer) {
	taskString := `normal task
sub folder task :subTaskFile sub_task.txt
task id 3
`

	taskBuf := bytes.NewBufferString(taskString)
	scanner := bufio.NewScanner(taskBuf)
	tasks, maxTaskID := task.CreateTasks(scanner, 0)

	config, buf := util.ReadTestConfigRelativePath("../")
	config.Task.TaskFolder = "../"

	s := &State{
		Tasks:     tasks,
		Config:    config,
		MaxTaskID: maxTaskID,
	}

	return s, buf
}

func TestReadSubTasks(t *testing.T) {
	cmd := NewReload()

	Convey("correct", t, func() {
		s, buf := createSubTaskTestState()

		Convey("should not change top level tasks", func() {
			So(len(s.Tasks), ShouldEqual, 3)
			So(s.MaxTaskID, ShouldEqual, 3)
		})

		Convey("should set sub task from file", func() {
			buf.Reset()
			cmd.readSubTaskFile(s)
			testTask := s.Tasks[1]

			Convey("output", func() {
				So(buf.String(), ShouldEqual, "read 3 tasks from ../sub_task.txt\n")
			})

			Convey("num", func() {
				So(len(testTask.SubTasks), ShouldEqual, 2)
				So(len(testTask.SubTasks[0].SubTasks), ShouldEqual, 1)
				So(len(testTask.SubTasks[1].SubTasks), ShouldEqual, 0)
			})

			Convey("title", func() {
				So(testTask.SubTasks[0].Name, ShouldEqual, "this is sub files task 1")
				So(testTask.SubTasks[0].SubTasks[0].Name, ShouldEqual, "this is sub files task sub task")
				So(testTask.SubTasks[1].Name, ShouldEqual, "this is sub files task 2")
			})

			Convey("id", func() {
				So(s.MaxTaskID, ShouldEqual, 6)
				So(testTask.SubTasks[0].ID, ShouldEqual, 4)
				So(testTask.SubTasks[0].SubTasks[0].ID, ShouldEqual, 5)
				So(testTask.SubTasks[1].ID, ShouldEqual, 6)
			})

			Convey("level", func() {
				level := testTask.Level
				So(testTask.SubTasks[0].Level, ShouldEqual, level+1)
				So(testTask.SubTasks[0].SubTasks[0].Level, ShouldEqual, level+2)
				So(testTask.SubTasks[1].Level, ShouldEqual, level+1)
			})
		})
	})

	Convey("invalid", t, func() {
		Convey("file not exist", func() {
			s, buf := createSubTaskTestState()
			beforeMaxID := s.MaxTaskID

			buf.Reset()
			testTask := s.Tasks[1]
			testTask.Attributes["subTaskFile"] = "test.txt"
			cmd.readSubTaskFile(s)
			So(buf.String(), ShouldEqual, "open ../test.txt: no such file or directory\n")
			So(s.MaxTaskID, ShouldEqual, beforeMaxID)
		})
	})
}
