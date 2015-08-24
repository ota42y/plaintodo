package output

import (
	"../ls"
	"io"
)

func outputShowTask(w io.Writer, task *ls.ShowTask, showID bool, taskLevel int) {
	if task.Task != nil {
		w.Write([]byte(task.Task.StringWithTaskLevel(showID, taskLevel)))
		w.Write([]byte("\n"))
	}

	for _, subTask := range task.SubTasks {
		outputShowTask(w, subTask, showID, taskLevel+1)
	}
}

// Output write task and sub tasks to io.Writer
func Output(w io.Writer, tasks []*ls.ShowTask, showID bool, baseTaskLevel int) {
	for _, task := range tasks {
		outputShowTask(w, task, showID, baseTaskLevel)
		w.Write([]byte("\n")) // make blank line between top level task
	}
	return
}
