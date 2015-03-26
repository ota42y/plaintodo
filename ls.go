package main

type ShowTask struct {
	Task     *Task
	SubTasks []*Task
}

func Ls(tasks []*Task) []*ShowTask {
	return make([]*ShowTask, 0)
}

func filterRoot(tasks []*Task) []*ShowTask {
	showTasks := make([]*ShowTask, 0)

	for _, task := range tasks{
		showTask := filter(task)
		if showTask != nil{
			showTasks = append(showTasks, showTask)
		}
	}

	return showTasks
}

func filter(task *Task) *ShowTask {
	return nil
}
