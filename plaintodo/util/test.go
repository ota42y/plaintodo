package util

import (
	"../config"
	"../task"
)

// ReadTestTasks read 'test_task.txt
func ReadTestTasks() []*task.Task {
	filename := "test_task.txt"
	tasks, _ := task.ReadTasks(filename)
	return tasks
}

// ReadTestConfig read 'test_config.toml'
func ReadTestConfig() *config.Config {
	return config.ReadConfig("test_config.toml")
}
