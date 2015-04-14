package main

import (
	   "os"
)

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
