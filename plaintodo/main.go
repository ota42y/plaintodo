package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

type Config struct {
	Paths PathConfig
}

type PathConfig struct {
	Task string
}

func readConfig(tomlFilepath string) *Config {
	var config Config
	_, err := toml.DecodeFile(tomlFilepath, &config)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &config
}

func main() {
	config := readConfig("config.toml")
	if config != nil {
		tasks := ReadTasks(config.Paths.Task)
		showTasks := Ls(tasks, NewExpireDateQuery("due", time.Now(), make([]Query, 0), make([]Query, 0)))
		Output(os.Stdout, showTasks)
	}
}
