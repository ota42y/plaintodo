package main

import (
	"testing"
)

func ReadTestTasks() []*Task {
	filename := "test_task.txt"
	return ReadTasks(filename)
}


func TestReadConfig(t *testing.T) {
	config := readConfig("test.toml")

	taskPath := "./test_task.txt"
	if config.TaskFilepath != taskPath {
		t.Errorf("config.TaskFilepath shuld be %s, but %s", taskPath, config.TaskFilepath)
		t.FailNow()
	}
}
