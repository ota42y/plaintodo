package util

import (
	"path"

	"../config"
	"../task"
)

// ReadTestTasks read 'test_task.txt
func ReadTestTasks() []*task.Task {
	filename := "test_task.txt"
	tasks, _ := task.ReadTasks(filename)
	return tasks
}

// ReadTestTaskRelativePath read 'test_task.txt' by relative path
func ReadTestTaskRelativePath(folder string) []*task.Task {
	tasks, _ := task.ReadTasks(path.Join(folder, "test_task.txt"))
	return tasks
}

// ReadTestConfig read 'test_config.toml'
func ReadTestConfig() *config.Config {
	return config.ReadConfig("test_config.toml")
}

// ReadTestConfigRelativePath read 'test_config.toml' by relative path
func ReadTestConfigRelativePath(folder string) *config.Config {
	return config.ReadConfig(path.Join(folder, "test_config.toml"))
}
