package main

import (
	"testing"
	"time"
)

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

func TestShowSubTasks(t *testing.T) {
	tasks := ReadTestTasks()

	query, _ := getQuery(" :id 2")
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

	if showTasks[0].SubTasks[0].SubTasks[0].Task.Id != 3 {
		t.Errorf("shuld return :id 3 task, but %d", showTasks[0].SubTasks[0].SubTasks[0].Task.Id)
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

	query, _ = getQuery(" :id 2")
	base = query.(*QueryBase)

	if len(base.and) == 0 {
		t.Errorf("IdQuery dosen't created")
		t.FailNow()
	}

	showTasks = ExecuteQuery(" :id 2", tasks)
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

	showTasks = ExecuteQuery(" :id 2 :no-sub-tasks", tasks)
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

	task = showTasks[0].SubTasks[0]
	if task.Task.Id != 2 {
		t.Errorf("Task.Id shuld be 2 but %d", task.Task.Id)
		t.FailNow()
	}

	if len(task.SubTasks) != 0 {
		t.Errorf("shuld no subtasks, but %d subtasks", len(task.SubTasks))
		t.FailNow()
	}

	cmd := NewCompleteCommand()
	cmd.completeTask(2, tasks)

	showTasks = ExecuteQuery(" :id 1 :complete", tasks)
	if len(showTasks) != 1 {
		t.Errorf("shuld return only one task, but %d tasks", len(showTasks))
		t.FailNow()
	}

	if showTasks[0].Task.Id != 1 {
		t.Errorf("shuld return task id = 1, but $v task", showTasks[0].Task)
		t.FailNow()
	}

	if len(showTasks[0].SubTasks) != 2 {
		t.Errorf("When :complete option set, get all completed sub tasks (2), but %d sub tasks", len(showTasks[0].SubTasks))
		t.FailNow()
	}

	task = showTasks[0].SubTasks[0]
	if task.Task.Id != 2 {
		t.Errorf("Task.Id shuld be 2 but %d", task.Task.Id)
		t.FailNow()
	}

	showTasks = ExecuteQuery(" :id 1", tasks)
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

	task = showTasks[0].SubTasks[0]
	if task.Task.Id != 4 {
		t.Errorf("Task.Id shuld be 4 but %d", task.Task.Id)
		t.FailNow()
	}

	if len(task.SubTasks) != 3 {
		t.Errorf("When :complete option not set, get all completed sub tasks (3), but %d sub tasks", len(task.SubTasks))
		t.FailNow()
	}

	// When no option set, get all started and no completed tasks
	tasks[0].SubTasks[1].SubTasks[1].Attributes["start"] = time.Now().Format(dateFormat)
	showTasks = ExecuteQuery("", tasks)
	if len(showTasks) != 2 {
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

	task = showTasks[0].SubTasks[0]
	if task.Task.Id != 4 {
		t.Errorf("Task.Id shuld be 4 but %d", task.Task.Id)
		t.FailNow()
	}

	if len(task.SubTasks) != 1 {
		t.Errorf("This shud be one sub task, but %d sub tasks", len(task.SubTasks))
		t.FailNow()
	}

	tasks = ReadTestTasks()
	cmd.completeTask(8, tasks)
	showTasks = ExecuteQuery("", tasks)
	if len(showTasks) != 1 {
		t.Errorf("if top level task completed, not show task, but %d task showed", len(showTasks))
		t.FailNow()
	}

	tasks = ReadTestTasks()
	delete(tasks[0].Attributes, "start")
	cmd.completeTask(2, tasks)
	cmd.completeTask(4, tasks)
	showTasks = ExecuteQuery("", tasks)
	if len(showTasks) != 1 {
		t.Errorf("if specific task is completed, don't show all parent")
		t.FailNow()
	}

	tasks = ReadTestTasks()
	showTasks = ExecuteQuery(" :overdue 2015-02-02", tasks)
	if len(showTasks) != 2 {
		t.Errorf("When start option set, return 2 tasks, but %d", len(showTasks))
		t.FailNow()
	}

	taskData := tasks[1].SubTasks[0]
	postpone := NewPostponeCommand()
	op := make(map[string]string)
	op["postpone"] = "3 day"
	postpone.postpone(taskData, op)

	showTasks = ExecuteQuery(" :overdue 2015-02-02", tasks)
	if len(showTasks) != 1 {
		t.Errorf("when postpone task, task isn't overdue but return %v", showTasks[1].Task)
		t.FailNow()
	}

	showTasks = ExecuteQuery(" :overdue 2015-02-05", tasks)
	if len(showTasks) != 2 {
		t.Errorf("when task ovordue in postpone time but return %d", len(showTasks))
		t.FailNow()
	}
}
