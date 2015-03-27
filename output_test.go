package main

import (
	"testing"
	"bytes"
)

func TestPrint(t *testing.T) {
	tasks := ReadTestTasks()
	showTasks := Ls(tasks)

	// show first task
	outputTasks := make([]*ShowTask, 1)
	outputTasks[0] = showTasks[0]

	buf := &bytes.Buffer{}
	Output(buf, outputTasks)

	
	correctString := `
go to SSA :due 2015-02-01
  create a set list :due 2015-01-31 :important
    add music to player
  buy items
    buy battery
    buy ultra orange
    buy king blade
`
	result := buf.String()
	if result != correctString {
		t.Errorf("incorrect string returned: %s", result)
		t.FailNow()
	}
}
