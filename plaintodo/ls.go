package main

import (
	"strconv"
)

type ShowTask struct {
	Task     *Task
	SubTasks []*ShowTask
}

func Ls(tasks []*Task, query Query) []*ShowTask {
	return filterRoot(tasks, query)
}

func filterRoot(tasks []*Task, query Query) []*ShowTask {
	showTasks := make([]*ShowTask, 0)

	for _, task := range tasks {
		showTask := filter(task, query)
		if showTask != nil {
			showTasks = append(showTasks, showTask)
		}
	}

	return showTasks
}

func filter(task *Task, query Query) *ShowTask {
	subTasks := make([]*ShowTask, 0)
	for _, task := range task.SubTasks {
		subTask := filter(task, query)
		if subTask != nil {
			subTasks = append(subTasks, subTask)
		}
	}

	// if SubTask exist, or query correct show parent task
	is_show := true
	if query != nil {
		is_show = len(subTasks) != 0 || query.Check(task)
	}

	if is_show {
		return &ShowTask{
			Task:     task,
			SubTasks: subTasks,
		}
	}

	return nil
}

func GetQuery(queryString string) Query {
	queryArray := ParseOptions(queryString)
	parent := NewQueryBase(make([]Query, 0), make([]Query, 0))

	for key, value := range queryArray {
		switch {
		case key == "level":
			{
				num, err := strconv.Atoi(value)
				if err == nil {
					parent.and = append(parent.and, NewMaxLevelQuery(num, make([]Query, 0), make([]Query, 0)))
				}
			}
		}
	}

	return parent
}
