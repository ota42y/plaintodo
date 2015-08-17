package ls

import (
	"testing"

	"../util"
)

func TestLs(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("../")

	showTasks := Ls(tasks, nil)
	if len(showTasks) != 2 {
		t.Errorf("top level shud be 2 but %d", len(showTasks))
		t.FailNow()
	}

	subTasks := showTasks[0].SubTasks
	if len(subTasks) != 2 {
		t.Errorf("top level shud be 2 but %d", len(subTasks))
		t.FailNow()
	}

	subTasks = subTasks[1].SubTasks
	if len(subTasks) != 3 {
		t.Errorf("top level shud be 3 but %d", len(subTasks))
		t.FailNow()
	}
}

func TestShowSubTasks(t *testing.T) {
	tasks := util.ReadTestTaskRelativePath("../")

	query, _ := GetQuery(" :id 2")
	showTasks := Ls(tasks, query)

	if len(showTasks) != 1 {
		t.Errorf("there is no task")
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("there is no sub task")
		t.FailNow()
	}

	if len(showTasks[0].SubTasks[0].SubTasks) != 0 {
		t.Errorf("there is sub task, but it shuld be 0")
		t.FailNow()
	}

	ShowAllChildSubTasks(showTasks)
	if len(showTasks) == 0 {
		t.Errorf("there is no show tasks")
		t.FailNow()
	}

	if len(showTasks) != 1 {
		t.Errorf("there is no task")
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("there is no sub task")
		t.FailNow()
	}

	if len(showTasks[0].SubTasks[0].SubTasks) != 1 {
		t.Errorf("there is no sub task")
		t.FailNow()
	}

	if len(showTasks[0].SubTasks[0].SubTasks[0].SubTasks) != 0 {
		t.Errorf("shuld return no sub tasks, but return %d", len(showTasks[0].SubTasks[0].SubTasks[0].SubTasks))
		t.FailNow()
	}

	if showTasks[0].SubTasks[0].SubTasks[0].Task.ID != 3 {
		t.Errorf("shuld return :id 3 task, but %d", showTasks[0].SubTasks[0].SubTasks[0].Task.ID)
		t.FailNow()
	}

}
