package output

import (
	"../ls"
	"io"
)

func outputShowTask(w io.Writer, task *ls.ShowTask, showID bool, taskLevel int, omit []string) {
	if task.Task != nil {
		w.Write([]byte(task.Task.StringWithTaskLevelAndOmit(showID, taskLevel, omit)))
		w.Write([]byte("\n"))
	}

	for _, subTask := range task.SubTasks {
		outputShowTask(w, subTask, showID, taskLevel+1, omit)
	}
}

// Output write task and sub tasks to io.Writer
func Output(w io.Writer, tasks []*ls.ShowTask, showID bool, baseTaskLevel int, omit []string) {
	for _, task := range tasks {
		outputShowTask(w, task, showID, baseTaskLevel, omit)
		w.Write([]byte("\n")) // make blank line between top level task
	}
	return
}
