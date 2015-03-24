package main

type ShowTask struct {
	Task     *Task
	SubTasks []*Task
}

func Ls(tasks []*Task) []*ShowTask {
	return make([]*ShowTask, 0)
}
