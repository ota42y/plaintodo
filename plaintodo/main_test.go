package main

import (
	"testing"
)

func ReadTestTasks() []*Task {
	filename := "test_task.txt"
	return ReadTasks(filename)
}

func TestReadConfig(t *testing.T) {
	config := readConfig("test_config.toml")

	taskPath := "./test_task.txt"
	if config.Paths.Task != taskPath {
		t.Errorf("config.Paths.Task shuld be %s, but %s", taskPath, config.Paths.Task)
		t.FailNow()
	}
}
