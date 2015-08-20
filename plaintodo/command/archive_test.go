package command

import (
	"testing"
	"time"

	"../util"
)

func TestGetCompleteDayList(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("..")
	cmd := NewArchive()

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
