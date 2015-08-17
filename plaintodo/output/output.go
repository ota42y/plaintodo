package output

import (
	"../ls"
	"io"
)

func outputShowTask(w io.Writer, task *ls.ShowTask, showID bool) {
	if task.Task != nil {
		w.Write([]byte(task.Task.String(showID)))
		w.Write([]byte("\n"))
	}

	for _, subTask := range task.SubTasks {
		outputShowTask(w, subTask, showID)
	}
}

// Output write task and sub tasks to io.Writer
func Output(w io.Writer, tasks []*ls.ShowTask, showID bool) {
	for _, task := range tasks {
		outputShowTask(w, task, showID)
		w.Write([]byte("\n")) // make blank line between top level task
	}
	return
}
