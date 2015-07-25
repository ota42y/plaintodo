package main

import (
	"testing"

	"./task"
)

func ReadTestTasks() []*task.Task {
	filename := "test_task.txt"
	tasks, _ := task.ReadTasks(filename)
	return tasks
}

func ReadTestConfig() *Config {
	return readConfig("test_config.toml")
}

func TestReadConfig(t *testing.T) {
	config := ReadTestConfig()

	taskPath := "./test_task.txt"
	if config.Paths.Task != taskPath {
		t.Errorf("config.Paths.Task shuld be %s, but %s", taskPath, config.Paths.Task)
		t.FailNow()
	}
}

func TestReadConfigError(t *testing.T) {
	config := readConfig("nothing")

	if config != nil {
		t.Errorf("if no file exist, return nil but return %v", config)
		t.FailNow()
	}
}
