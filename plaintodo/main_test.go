package main

import (
	"testing"

	"./config"
	"./util"
)

func TestReadConfig(t *testing.T) {
	config := util.ReadTestConfig()

	taskPath := "./test_task.txt"
	if config.Paths.Task != taskPath {
		t.Errorf("config.Paths.Task shuld be %s, but %s", taskPath, config.Paths.Task)
		t.FailNow()
	}
}

func TestReadConfigError(t *testing.T) {
	config := config.ReadConfig("nothing")

	if config != nil {
		t.Errorf("if no file exist, return nil but return %v", config)
		t.FailNow()
	}
}
