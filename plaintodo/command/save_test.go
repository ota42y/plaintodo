package command

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"../task"
	"../util"
)

func TestSave(t *testing.T) {
	cmd := NewSave()

	Convey("correct", t, func() {
		Convey("save tasks", func() {
			c, _ := util.ReadTestConfigRelativePath("..")
			c.Task.TaskFolder = "../../result"
			c.Task.DefaultFilename = "savetest.txt"

			s := &State{
				Config: c,
			}

			s.Tasks = make([]*task.Task, 0)
			cmd.saveToFile(s)

			s.Tasks = util.ReadTestTaskRelativePath("..")
			s.Tasks[0].SubTasks = make([]*task.Task, 0)
			cmd.saveToFile(s)

			loadTasks, _, _ := task.ReadTasks("../../result/savetest.txt", 0)

			Convey("top level tasks", func() {
				So(len(loadTasks), ShouldEqual, len(s.Tasks))
			})

			Convey("sub tasks", func() {
				So(len(loadTasks[0].SubTasks), ShouldEqual, len(s.Tasks[0].SubTasks))
			})
		})

		Convey("getSaveTaskList", func() {
			s, _ := createSubTaskTestState()
			reload := NewReload()
			reload.readSubTaskFile(s)

			list := cmd.getSaveTaskList(s)

			Convey("2 objects", func() {
				So(len(*list), ShouldEqual, 2)
			})

			Convey("correct filename", func() {
				So((*list)[0].FileName, ShouldEqual, "sub_task.txt")
				So(len((*list)[0].Tasks), ShouldEqual, 2)
				So((*list)[1].FileName, ShouldEqual, "test_task.txt")
				So(len((*list)[1].Tasks), ShouldEqual, 3)
			})
		})
	})
}
