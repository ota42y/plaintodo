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

	config := readConfig("config.toml")
	config.Writer = os.Stdout
	if config != nil {
		l := NewLiner(config, cmds)
		l.automaton.Execute("reload")
		l.Start()
	}
}
