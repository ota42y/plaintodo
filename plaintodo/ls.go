package main

import (
	"strconv"
	"time"
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

// show all sub task in selected task
func ShowAllChildSubTasks(showTasks []*ShowTask) {
	for _, task := range showTasks {
		showSubTasks(task)
	}
}

func showSubTasks(task *ShowTask) {
	if len(task.SubTasks) == 0 {
		// overwrite all sub tasks
		task.SubTasks = filterRoot(task.Task.SubTasks, nil)
		return
	}

	for _, subTask := range task.SubTasks {
		showSubTasks(subTask)
	}
}

func getQuery(queryString string) (query Query, queryMap map[string]string) {
	queryMap = ParseOptions(queryString)
	parent := NewQueryBase(make([]Query, 0), make([]Query, 0))

	for key, value := range queryMap {
		switch {
		case key == "level":
			{
				num, err := strconv.Atoi(value)
				if err == nil {
					parent.and = append(parent.and, NewMaxLevelQuery(num, make([]Query, 0), make([]Query, 0)))
				}
			}

		case key == "id":
			{
				num, err := strconv.Atoi(value)
				if err == nil {
					parent.and = append(parent.and, NewIdQuery(num, make([]Query, 0), make([]Query, 0)))
				}
			}
		}
	}

	return parent, queryMap
}

func ExecuteQuery(queryString string, tasks []*Task) []*ShowTask{
	var query Query = nil
	queryMap := make(map[string]string)
	if queryString != "" {
		// GetCommand expected ' :key value :key value', but option give ':key value :key value'
		// so add space to first
		query, queryMap = getQuery(" " + queryString)
	} else {
		query = NewBeforeDateQuery("due", time.Now(), make([]Query, 0), make([]Query, 0))
	}

	showTasks := Ls(tasks, query)
	_, ok := queryMap["subtask"]
	if ok {
		ShowAllChildSubTasks(showTasks)
	}

	return showTasks
}