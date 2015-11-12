package main

import (
	"fmt"
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
	cmds["task"] = command.NewAddTask()
	cmds["set"] = command.NewSetAttribute()
	cmds["start"] = command.NewStart()
	cmds["postpone"] = command.NewPostpone()
	cmds["move"] = NewMoveCommand()
	cmds["open"] = NewOpenCommand()
	cmds["nice"] = command.NewNice()
	cmds["alias"] = NewAliasCommand()

	c := config.ReadConfig("config.toml")
	if c != nil {
		c.Writer = os.Stdout
		l := NewLiner(c, cmds)
		l.e.Execute("reload")

		if len(os.Args) == 1 {
			l.Start()
		} else {
			l.e.Execute(fmt.Sprintf("%s", os.Args[1]))
		}
	}
}
