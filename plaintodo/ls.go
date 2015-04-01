package main

type ShowTask struct {
	Task     *Task
	SubTasks []*ShowTask
}

type Query interface {
	Check(task *Task) bool
}

type KeyValueQuery struct {
	key   string
	value string

	and []Query
	or  []Query
}

func (query KeyValueQuery) Check(task *Task) bool {
	if task == nil {
		return false
	}

	if task.Attributes[query.key] == query.value {
		for _, q := range query.and {
			if !q.Check(task) {
				return false
			}
		}
		return true
	} else {
		for _, q := range query.or {
			if q.Check(task) {
				return true
			}
		}
		return false
	}
}

func NewKeyValueQuery(key string, value string, and []Query, or []Query) *KeyValueQuery {
	return &KeyValueQuery{
		key:   key,
		value: value,
		and:   and,
		or:    or,
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