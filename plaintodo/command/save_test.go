package command

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"time"

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

func TestGetCompleteDayList(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("..")
	cmd := NewSave()

	testTimeList := [...]string{"2015-01-31 10:42", "2015-01-29", "2015-01-30 10:42", "2015-01-30"}
	tasks[0].Attributes["complete"] = testTimeList[0]
	tasks[0].SubTasks[0].Attributes["complete"] = testTimeList[1]
	tasks[0].SubTasks[1].Attributes["complete"] = testTimeList[2]
	tasks[0].SubTasks[1].SubTasks[0].Attributes["complete"] = testTimeList[3]

	correctTimeList := make([]time.Time, 3)
	parseList := [...]int{1, 2, 0}
	for index, value := range parseList {
		timeData, ok := util.ParseTime(testTimeList[value])
		if !ok {
			t.Errorf("parse error %s", testTimeList[value])
			t.FailNow()
		}
		correctTimeList[index] = timeData
	}

	timeList := cmd.getCompleteDayList(tasks)
	if len(timeList) != len(correctTimeList) {
		t.Errorf("shuld return %d items, but %d items %v", len(correctTimeList), len(timeList), timeList)
		t.FailNow()
	}

	for index, item := range correctTimeList {
		year, month, day := item.Date()
		y, m, d := timeList[index].Date()

		if (year != y) || (month != m) || (day != d) {
			t.Errorf("shuld return %v but %v", item, timeList[index])
			t.FailNow()
		}
	}
}
