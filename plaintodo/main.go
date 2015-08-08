package main

import (
	"os"

	"./config"
	"./executor"
)

func main() {
	cmds := make(map[string]executor.Command)
	cmds["exit"] = NewExitCommand()
	cmds["reload"] = NewReloadCommand()
	cmds["ls"] = NewLsCommand(os.Stdout)
	cmds["lsall"] = NewLsAllCommand(os.Stdout)
	cmds["save"] = NewSaveCommand()
	cmds["complete"] = NewCompleteCommand()
	cmds["task"] = NewAddTaskCommand()
	cmds["subtask"] = NewAddSubTaskCommand()
	cmds["set"] = NewSetAttributeCommand()
	cmds["start"] = NewStartCommand()
	cmds["postpone"] = NewPostponeCommand()
	cmds["move"] = NewMoveCommand()
	cmds["open"] = NewOpenCommand()
	cmds["nice"] = NewNiceCommand()
	cmds["alias"] = NewAliasCommand()

	c := config.ReadConfig("config.toml")
	c.Writer = os.Stdout
	if c != nil {
		l := NewLiner(c, cmds)
		l.e.Execute("reload")
		l.Start()
	}
}
