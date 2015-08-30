package main

import (
	"os"

	"./command"
	"./config"
)

func main() {
	cmds := make(map[string]command.Command)
	cmds["exit"] = command.NewExit()
	cmds["reload"] = command.NewReload()
	cmds["ls"] = NewLsCommand(os.Stdout)
	cmds["lsall"] = NewLsAllCommand(os.Stdout)
	cmds["update"] = command.NewUpdate()
	//cmds["save"] = command.NewSave()
	//cmds["archive"] = command.NewArchive()
	cmds["complete"] = command.NewComplete()
	cmds["task"] = NewAddTaskCommand()
	cmds["subtask"] = NewAddSubTaskCommand()
	cmds["set"] = command.NewSetAttribute()
	cmds["start"] = NewStartCommand()
	cmds["postpone"] = command.NewPostpone()
	cmds["move"] = NewMoveCommand()
	cmds["open"] = NewOpenCommand()
	cmds["nice"] = command.NewNice()
	cmds["alias"] = NewAliasCommand()

	c := config.ReadConfig("config.toml")
	c.Writer = os.Stdout
	if c != nil {
		l := NewLiner(c, cmds)
		l.e.Execute("reload")
		l.Start()
	}
}
