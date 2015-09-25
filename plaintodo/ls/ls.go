package ls

import (
	"strconv"
	"time"

	"../query"
	"../task"
	"../util"
)

// ShowTask is result of Ls function
// If query select sub task, it's in the ShowTask.SubTasks
// If query not select sub task, it's not in that but Task.SubTasks contain it.
type ShowTask struct {
	Task     *task.Task // specific task
	SubTasks []*ShowTask
}

// Ls return tasks which filter by query
func Ls(tasks []*task.Task, q query.Query) []*ShowTask {
	return filterRoot(tasks, q)
}

func filterRoot(tasks []*task.Task, q query.Query) []*ShowTask {
	var showTasks []*ShowTask
	for _, task := range tasks {
		showTask := Filter(task, q)
		if showTask != nil {
			showTasks = append(showTasks, showTask)
		}
	}

	return showTasks
}

// Filter filtering tasks by query
// This should be private
func Filter(task *task.Task, q query.Query) *ShowTask {
	var subTasks []*ShowTask
	for _, task := range task.SubTasks {
		subTask := Filter(task, q)
		if subTask != nil {
			subTasks = append(subTasks, subTask)
		}
	}

	// if SubTask exist, or query correct show parent task
	isShow := true
	if q != nil {
		isShow = len(subTasks) != 0 || q.Check(task)
	}

	if isShow {
		return &ShowTask{
			Task:     task,
			SubTasks: subTasks,
		}
	}

	return nil
}

// DeleteAllCompletedTasks delete all task which have complete attribute
func DeleteAllCompletedTasks(showTasks []*ShowTask) []*ShowTask {
	var newSubTasks []*ShowTask
	for _, task := range showTasks {
		if deleteAllCompletedSubTasks(task) {
			newSubTasks = append(newSubTasks, task)
		}
	}
	return newSubTasks
}
func deleteAllCompletedSubTasks(task *ShowTask) bool {
	// check all sub tasks
	var newSubTasks []*ShowTask
	for _, subTask := range task.SubTasks {
		if deleteAllCompletedSubTasks(subTask) {
			newSubTasks = append(newSubTasks, subTask)
		}
	}

	task.SubTasks = newSubTasks
	if len(newSubTasks) == 0 {
		// no sub task or all sub task is completed

		// if not show complete task
		_, isComplete := task.Task.Attributes["complete"]

		if isComplete {
			return false
		}
	}

	// task exist or no completed task
	return true
}

// ShowAllChildSubTasks show all sub task in selected task
func ShowAllChildSubTasks(showTasks []*ShowTask) {
	for _, task := range showTasks {
		showSubTasks(task)
	}
}

func showSubTasks(task *ShowTask) {
	if len(task.SubTasks) == 0 {
		// overwrite all sub tasks
		task.SubTasks = filterRoot(task.Task.SubTasks, nil)
	}

	// check all sub tasks
	var newSubTasks []*ShowTask
	for _, subTask := range task.SubTasks {
		showSubTasks(subTask)
		newSubTasks = append(newSubTasks, subTask)
	}

	task.SubTasks = newSubTasks
	return
}

// GetQuery parse string to query
func GetQuery(queryString string) (query.Query, map[string]string) {
	queryMap := task.ParseOptions(queryString)
	parent := query.NewBase(make([]query.Query, 0), make([]query.Query, 0))

	for key, value := range queryMap {
		switch {
		case key == "level":
			{
				num, err := strconv.Atoi(value)
				if err == nil {
					parent.And = append(parent.And, query.NewMaxLevel(num, make([]query.Query, 0), make([]query.Query, 0)))
				}
			}

		case key == "id":
			{
				num, err := strconv.Atoi(value)
				if err == nil {
					parent.And = append(parent.And, query.NewID(num, make([]query.Query, 0), make([]query.Query, 0)))
				}
			}

		case key == "overdue":
			{
				t, ok := util.ParseTime(value)
				if ok {
					noPostpone := query.NewNoKey("postpone", make([]query.Query, 0), make([]query.Query, 0))
					noPostpone.And = append(noPostpone.And, query.NewBeforeDate("start", t, make([]query.Query, 0), make([]query.Query, 0)))

					overduePostpone := query.NewBeforeDate("postpone", t, make([]query.Query, 0), make([]query.Query, 0))
					overduePostpone.Or = append(overduePostpone.Or, noPostpone)

					parent.And = append(parent.And, overduePostpone)
				}
			}
		}
	}

	// if not complete show
	_, isCompleteShow := queryMap["complete"]
	if !isCompleteShow {
		parent.And = append(parent.And, query.NewNoKey("complete", make([]query.Query, 0), make([]query.Query, 0)))
	}

	return parent, queryMap
}

// ExecuteQuery execute query and return tasks
func ExecuteQuery(queryString string, tasks []*task.Task) ([]*ShowTask, bool) {
	if queryString == "" {
		// default query
		queryString = ":omit :overdue " + time.Now().Format(util.DateTimeFormat)
	}

	// GetCommand expected ' :key value :key value', but option give ':key value :key value'
	// so add space to first
	query, queryMap := GetQuery(" " + queryString)
	showTasks := Ls(tasks, query)

	_, ok := queryMap["no-sub-tasks"]
	if !ok {
		ShowAllChildSubTasks(showTasks)
	}

	_, isOmit := queryMap["omit"]

	// if not complete query, show only no completed query
	_, isCompleteShow := queryMap["complete"]
	if !isCompleteShow {
		showTasks = DeleteAllCompletedTasks(showTasks)
	}

	return showTasks, isOmit
}
