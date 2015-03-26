package main

type ShowTask struct {
	Task     *Task
	SubTasks []*ShowTask
}

func Ls(tasks []*Task) []*ShowTask {
	return filterRoot(tasks)
}

func filterRoot(tasks []*Task) []*ShowTask {
	showTasks := make([]*ShowTask, 0)

	for _, task := range tasks {
		showTask := filter(task)
		if showTask != nil {
			showTasks = append(showTasks, showTask)
		}
	}

	return showTasks
}

func filter(task *Task) *ShowTask {
	subTasks := make([]*ShowTask, 0)
	for _, task := range task.SubTasks {
		subTask := filter(task)
		if subTask != nil {
			subTasks = append(subTasks, subTask)
		}
	}

	// if SubTask exist, show parent task
	//if len(subTasks) != 0 {
	// always return ShowTask
	return &ShowTask{
		Task:     task,
		SubTasks: subTasks,
	}
	//}

	//return nil
}
