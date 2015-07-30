package main

import (
	"fmt"
	"strings"

	"./task"
)

var optionSplit = " "

type Command interface {
	Execute(option string, automaton *Automaton) (terminate bool)
}

type Automaton struct {
	Tasks          []*task.Task
	MaxTaskID      int
	Commands       map[string]Command
	CommandAliases map[string]string
	Config         *Config // config.go
}

func NewAutomaton(config *Config, commands map[string]Command) *Automaton {
	aliases := make(map[string]string)

	if config != nil {
		for _, alias := range config.Command.Aliases {
			aliases[alias.Name] = alias.Command
		}
	}

	return &Automaton{
		Tasks:          make([]*task.Task, 0),
		Commands:       commands,
		CommandAliases: aliases,
		Config:         config,
	}
}

// cmd shuld be "cmd options"
func (a *Automaton) Execute(command string) (terminate bool) {
	splits := strings.SplitAfterN(command, optionSplit, 2)
	if len(splits) == 0 {
		// no command
		return false
	}

	cmd := strings.TrimSpace(splits[0])

	option := ""
	if 1 < len(splits) {
		option = strings.TrimSpace(splits[1])
	}

	alias, ok := a.CommandAliases[cmd]
	if ok {
		newCommand := alias + " " + option
		fmt.Fprintf(a.Config.Writer, "alias %s = %s\n", cmd, alias)
		fmt.Fprintf(a.Config.Writer, "command: %s\n", newCommand)

		return a.Execute(newCommand)
	}

	value, ok := a.Commands[cmd]
	if a.Config != nil && a.Config.Writer != nil {
		if ok {
			a.Config.Writer.Write([]byte(cmd + " hit\n"))
		} else {
			a.Config.Writer.Write([]byte(cmd + " not hit\n"))
		}
	}

	if ok {
		return value.Execute(option, a)
	} else {
		// no command
		return false
	}
}
