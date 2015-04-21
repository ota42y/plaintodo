package main

import (
	"os"
)

func main() {
	cmds := make(map[string]Command)
	cmds["exit"] = NewExitCommand()
	cmds["reload"] = NewReloadCommand()
	cmds["ls"] = NewLsCommand(os.Stdout)
	cmds["lsall"] = NewLsAllCommand(os.Stdout)
	cmds["save"] = NewSaveCommand()

	config := readConfig("config.toml")
	config.Writer = os.Stdout
	if config != nil {
		l := NewLiner(config, cmds)
		l.Start()
	}
}
