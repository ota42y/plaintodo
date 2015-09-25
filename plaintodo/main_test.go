package main

import (
	"testing"

	"./config"
	"./util"
)

func TestReadConfig(t *testing.T) {
	config, _ := util.ReadTestConfig()

	taskPath := "test_task.txt"
	if config.Task.DefaultFilename != taskPath {
		t.Errorf("config.Paths.Task shuld be %s, but %s", taskPath, config.Task.DefaultFilename)
		t.FailNow()
	}

	if len(config.Command.Omits) == 0 {
		t.Errorf("config.Command.Omits not loaded")
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
