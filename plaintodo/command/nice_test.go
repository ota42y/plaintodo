package command

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"../task"
	"../util"
)

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
			So(buf.String(), ShouldEqual, "Done nice\nevernote url change 1 tasks\n")
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
			So(buf.String(), ShouldEqual, "Done nice\nevernote url change 1 tasks\n")
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
