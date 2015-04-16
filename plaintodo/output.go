package main

import (
	"io"
)

func outputShowTask(w io.Writer, task *ShowTask) {
	if task.Task != nil {
		w.Write([]byte(task.Task.String()))
		w.Write([]byte("\n"))
	}

	for _, subTask := range task.SubTasks {
		outputShowTask(w, subTask)
	}
}

func Output(w io.Writer, tasks []*ShowTask) {
	for _, task := range tasks {
		outputShowTask(w, task)
		w.Write([]byte("\n")) // make blank line between top level task
	}
	return
}
