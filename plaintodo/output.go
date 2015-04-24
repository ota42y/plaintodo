package main

import (
	"io"
)

func outputShowTask(w io.Writer, task *ShowTask, showId bool) {
	if task.Task != nil {
		w.Write([]byte(task.Task.String(showId)))
		w.Write([]byte("\n"))
	}

	for _, subTask := range task.SubTasks {
		outputShowTask(w, subTask, showId)
	}
}

func Output(w io.Writer, tasks []*ShowTask, showId bool) {
	for _, task := range tasks {
		outputShowTask(w, task, showId)
		w.Write([]byte("\n")) // make blank line between top level task
	}
	return
}
