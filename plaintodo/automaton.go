package main

import (
	"strings"
)

var optionSplit = " "

type Command interface {
	Execute(option string, automaton *Automaton) (terminate bool)
}

type Automaton struct {
	Tasks    []*Task
	Commands map[string]Command
	Config   *Config // main.go
}

func NewAutomaton(config *Config, commands map[string]Command) *Automaton {
	return &Automaton{
		Tasks:    make([]*Task, 0),
		Commands: commands,
		Config:   config,
	}
}

// cmd shuld be "cmd options"
func (a *Automaton) Execute(command string) (terminate bool) {
	splits := strings.SplitAfterN(command, optionSplit, 2)
	if len(splits) == 0{
		// no command
		return false
	}

	cmd := strings.TrimSpace(splits[0])

	option := ""
	if 1 < len(splits){
		option = strings.TrimSpace(splits[1])
	}

	value, ok := a.Commands[cmd]
	if ok {
		return value.Execute(option, a)
	} else {
		// no command
		return false
	}
}
