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

	query, _ := getQuery(" :level 2")
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

	query, _ = getQuery(" :subtask :id 2")
	base = query.(*QueryBase)

	if len(base.and) == 0 {
		t.Errorf("IdQuery dosen't created")
		t.FailNow()
	}

	showTasks = ExecuteQuery(" :subtask :id 2", tasks)
	if len(showTasks) != 1 {
		t.Errorf("shuld return only one task, but %d tasks", len(showTasks))
		t.FailNow()
	}

	if showTasks[0].Task.Id != 1 {
		t.Errorf("shuld return task id = 1, but $v task", showTasks[0].Task)
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 1 {
		t.Errorf("one sub task shuld be exsit, but %d task", len(showTasks[0].SubTasks))
		t.FailNow()
	}

	task := showTasks[0].SubTasks[0]
	if task.Task.Id != 2 {
		t.Errorf("Task.Id shuld be 2 but %d", task.Task.Id)
		t.FailNow()
	}

	if task.Task.Name != "create a set list" {
		t.Errorf("shuld return task id = 2, but $v task", task.Task)
		t.FailNow()
	}

	if len(task.SubTasks) != len(task.Task.SubTasks) {
		t.Errorf("When :subtask option set, get %d sub tasks, but %d sub tasks", len(task.Task.SubTasks), len(task.SubTasks))
		t.FailNow()
	}

	subTask := task.SubTasks[0]
	if subTask.Task.Id != 3 {
		t.Errorf("Get :id 3 task, but %d", subTask.Task.Id)
		t.FailNow()
	}

	if len(subTask.SubTasks) != len(subTask.Task.SubTasks) {
		t.Errorf("When :subtask option set, get %d sub tasks, but %d sub tasks", len(subTask.Task.SubTasks), len(subTask.SubTasks))
		t.FailNow()
	}

}
