package command

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"../task"
	"../util"
	"fmt"
	"time"
)

func execFixDateInKey(cmd *Nice, key string, taskString string, today time.Time) *task.Task {
	testTask, _ := task.NewTask(taskString, 1)
	cmd.fixDateInKey(testTask, key, today)

	return testTask
}

func TestReplaceDate(t *testing.T) {
	cmd := NewNice()
	today := time.Now()

	Convey("correct", t, func() {
		todayDateStr := today.Format(util.DateFormat)
		todayTimeStr := today.Format(util.DateTimeFormat)

		Convey("start", func() {
			key := "start"

			Convey("now", func() {
				value := "now"
				taskString := fmt.Sprintf("test task :%s %s", key, value)
				testTask := execFixDateInKey(cmd, key, taskString, today)
				So(testTask.Attributes[key], ShouldEqual, todayTimeStr)
			})

			Convey("today", func() {
				value := "today"
				taskString := fmt.Sprintf("test task :%s %s", key, value)
				testTask := execFixDateInKey(cmd, key, taskString, today)
				So(testTask.Attributes[key], ShouldEqual, todayDateStr)
			})

			Convey("today 12:00", func() {
				value := "today 12:00"
				taskString := fmt.Sprintf("test task :%s %s", key, value)
				testTask := execFixDateInKey(cmd, key, taskString, today)
				So(testTask.Attributes[key], ShouldEqual, todayDateStr+" 12:00")
			})
		})

		Convey("postpone", func() {
			key := "postpone"
			startData := "2015-02-01"

			Convey("now", func() {
				value := "now"
				taskString := fmt.Sprintf("test task :%s %s :start %s", key, value, startData)
				testTask := execFixDateInKey(cmd, key, taskString, today)
				So(testTask.Attributes[key], ShouldEqual, todayTimeStr)
				So(testTask.Attributes["start"], ShouldEqual, startData)
			})

			Convey("today", func() {
				value := "today"
				taskString := fmt.Sprintf("test task :%s %s :start %s", key, value, startData)
				testTask := execFixDateInKey(cmd, key, taskString, today)
				So(testTask.Attributes[key], ShouldEqual, todayDateStr)
				So(testTask.Attributes["start"], ShouldEqual, startData)
			})

			Convey("today 12:00", func() {
				value := "today 12:00"
				taskString := fmt.Sprintf("test task :%s %s :start %s", key, value, startData)
				testTask := execFixDateInKey(cmd, key, taskString, today)
				So(testTask.Attributes[key], ShouldEqual, todayDateStr+" 12:00")
				So(testTask.Attributes["start"], ShouldEqual, startData)
			})
		})
	})

	Convey("incorrect", t, func() {
		Convey("not change task", func() {
			taskString := "test task :postpone 2015-02-01 12:00 :start 2015-02-01 11:00"

			testTask, _ := task.NewTask(taskString, 1)
			cmd.fixDateInKey(testTask, "start", today)
			cmd.fixDateInKey(testTask, "postpone", today)

			So(testTask.Attributes["postpone"], ShouldEqual, "2015-02-01 12:00")
			So(testTask.Attributes["start"], ShouldEqual, "2015-02-01 11:00")
		})
	})
}

func TestChangeDate(t *testing.T) {
	cmd := NewNice()

	config, buf := util.ReadTestConfigRelativePath("..")
	s := &State{
		Config: config,
	}

	Convey("correct", t, func() {
		s.Tasks = util.ReadTestTasks()

		testTask, _ := task.NewTask("test task for today :start now", 1)
		s.Tasks = append(s.Tasks, testTask)

		testTask, _ = task.NewTask("test task for today :start today", 1)
		s.Tasks = append(s.Tasks, testTask)

		testTask, _ = task.NewTask("test task for today :start today 12:00", 1)
		s.Tasks = append(s.Tasks, testTask)

		testTask, _ = task.NewTask("test task for today :postpone now :start 2015-02-01", 1)
		s.Tasks = append(s.Tasks, testTask)

		testTask, _ = task.NewTask("test task for today :postpone today :start 2015-02-01", 1)
		s.Tasks = append(s.Tasks, testTask)

		testTask, _ = task.NewTask("test task for today :postpone today 12:00 :start 2015-02-01", 1)
		s.Tasks = append(s.Tasks, testTask)

		buf.Reset()
		terminate := cmd.Execute("", s)
		So(terminate, ShouldBeFalse)
		So(buf.String(), ShouldEqual, "Done nice\nevernote url change 0 tasks\nchange 6 tasks date\n")
	})

	Convey("no change tasks", t, func() {
		s.Tasks = util.ReadTestTasks()

		buf.Reset()
		terminate := cmd.Execute("", s)
		So(terminate, ShouldBeFalse)
		So(buf.String(), ShouldEqual, "Done nice\nevernote url change 0 tasks\nchange 0 tasks date\n")
	})
}

func TestNiceCommand(t *testing.T) {
	cmd := NewNice()

	config, buf := util.ReadTestConfigRelativePath("..")
	s := &State{
		Config: config,
	}

	evernoteURL := "https://www.evernote.com/shard/s1/nl/111111/abfdef-ght1234567890"
	correctURL := "evernote:///view/111111/s1/abfdef-ght1234567890/abfdef-ght1234567890/"

	Convey("correct", t, func() {
		Convey("there isn^t / in bottom", func() {
			testTask, _ := task.NewTask("task test :url "+evernoteURL, 1)
			s.Tasks = make([]*task.Task, 0)
			s.Tasks = append(s.Tasks, testTask)

			buf.Reset()
			terminate := cmd.Execute("", s)
			So(terminate, ShouldBeFalse)
			So(buf.String(), ShouldEqual, "Done nice\nevernote url change 1 tasks\nchange 0 tasks date\n")
			So(testTask.Attributes["url"], ShouldEqual, correctURL)
		})

		Convey("there is / in bottom", func() {
			slashURL := evernoteURL + "/"
			testTask, _ := task.NewTask("task test :url "+slashURL, 1)

			s.Tasks = make([]*task.Task, 0)
			s.Tasks = append(s.Tasks, testTask)

			buf.Reset()
			terminate := cmd.Execute("", s)
			So(terminate, ShouldBeFalse)
			So(buf.String(), ShouldEqual, "Done nice\nevernote url change 1 tasks\nchange 0 tasks date\n")
			So(testTask.Attributes["url"], ShouldEqual, correctURL)
		})

		Convey("not change other task", func() {
			s.Tasks = util.ReadTestTaskRelativePath("..")
			buf.Reset()
			terminate := cmd.Execute("", s)
			So(terminate, ShouldBeFalse)
			So(s.Tasks[1].SubTasks[0].Attributes["url"], ShouldEqual, "http://ota42y.com")
		})
	})
}
