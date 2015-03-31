package main

import (
	"bytes"
	"testing"
	"strings"
)

func TestPrint(t *testing.T) {
	tasks := ReadTestTasks()
	showTasks := Ls(tasks)

	// show first task
	outputTasks := make([]*ShowTask, 1)
	outputTasks[0] = showTasks[0]

	buf := &bytes.Buffer{}
	Output(buf, outputTasks)

	correctString := `go to SSA :due 2015-02-01
  create a set list :due 2015-01-31 :important
    add music to player
  buy items
    buy battery
    buy ultra orange
    buy king blade
`
	results := strings.Split(buf.String(), "\n")

	corrects := strings.Split(correctString, "\n")

	if len(corrects) != len(results) {
		t.Errorf("return %d strings bud %d", len(corrects), len(results))
		t.FailNow()
	}

	for index, str := range corrects{
		if results[index] != str {
			t.Errorf("return shuld be '%s', but '%s'", str, results[index])
			t.FailNow()
		}
	}



}
