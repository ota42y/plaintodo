package main

import "testing"

func TestLs(t *testing.T) {
	tasks := ReadTestTasks()

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

func TestGetQuery(t *testing.T) {
	tasks := ReadTestTasks()

	query := GetQuery(" :level 2")
	base := query.(*QueryBase)

	if len(base.and) == 0 {
		t.Errorf("MaxLevelQuery dosen't created")
		t.FailNow()
	}

	showTasks := Ls(tasks, query)

	if len(showTasks) == 0 {
		t.Errorf("return no tasks")
		t.FailNow()
	}

	for _, task := range showTasks {
		if 2 < task.Task.Level {
			t.Errorf("Task.Level shuld be less than 2 but %v", task.Task)
			t.FailNow()
		}
		for _, subTask := range task.SubTasks {
			if 2 < subTask.Task.Level {
				t.Errorf("Task.Level shuld be less than 2 but %v", task.SubTasks)
				t.FailNow()
			}
			if len(subTask.SubTasks) != 0 {
				t.Errorf("SubTasks shuld be 0 but %d", len(subTask.SubTasks))
				t.FailNow()
			}
		}
	}
}
