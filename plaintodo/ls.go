package main

type ShowTask struct {
	Task     *Task
	SubTasks []*ShowTask
}

type Query interface {
	Check(task *Task) bool
}

type QueryBase struct {
	and []Query
	or  []Query
}

func (query QueryBase) checkSubQuery(task *Task, isShow bool) bool {
	// If this query return true, check all and query
	// (Even if or query exist, we don't need check these.
	if isShow {
		for _, q := range query.and {
			if !q.Check(task) {
				return false
			}
		}
		return true
	} else {
		// If this query return false, check all or query
		// Even if and query exist, we don't need check these.
		for _, q := range query.or {
			if q.Check(task) {
				return true
			}
		}
		return false
	}
}

type KeyValueQuery struct {
	*QueryBase

	key   string
	value string
}

func (query KeyValueQuery) Check(task *Task) bool {
	if task == nil {
		return false
	}

	return query.checkSubQuery(task, task.Attributes[query.key] == query.value)
}

func NewKeyValueQuery(key string, value string, and []Query, or []Query) *KeyValueQuery {
	return &KeyValueQuery{
		QueryBase: &QueryBase{
			and: and,
			or:  or,
		},
		key:   key,
		value: value,
	}
}

func Ls(tasks []*Task) []*ShowTask {
	return filterRoot(tasks, nil)
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
