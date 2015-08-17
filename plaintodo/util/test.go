package util

import (
	"bytes"
	"path"

	"../config"
	"../task"
)

// ReadTestTasks read 'test_task.txt
func ReadTestTasks() []*task.Task {
	filename := "test_task.txt"
	tasks, _, _ := task.ReadTasks(filename, 0)
	return tasks
}

// ReadTestTaskRelativePath read 'test_task.txt' by relative path
func ReadTestTaskRelativePath(folder string) []*task.Task {
	tasks, _, _ := task.ReadTasks(path.Join(folder, "test_task.txt"), 0)
	return tasks
}

// ReadTestConfig read 'test_config.toml'
func ReadTestConfig() (*config.Config, *bytes.Buffer) {
	return ReadTestConfigRelativePath(".")
}

// ReadTestConfigRelativePath read 'test_config.toml' by relative path
func ReadTestConfigRelativePath(folder string) (*config.Config, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	c := config.ReadConfig(path.Join(folder, "test_config.toml"))
	c.Writer = buf
	return c, buf
}
