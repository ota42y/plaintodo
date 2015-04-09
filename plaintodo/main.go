package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
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
	cmds := make(map[string]Command)
	cmds["exit"] = NewExitCommand()
	cmds["reload"] = NewReloadCommand()
	cmds["ls"] = NewLsCommand(os.Stdout)

	config := readConfig("config.toml")
	if config != nil {
		l := NewLiner(config, cmds)
		l.Start()
	}
}
