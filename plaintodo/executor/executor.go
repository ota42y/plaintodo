package executor

import (
	"fmt"
	"strings"

	"../config"
	"../task"
)

var optionSplit = " "

// Command is command interface
type Command interface {
	Execute(option string, s *State) (terminate bool)
}

// Executor executor command
type Executor struct {
	S *State
}

// NewExecutor return Executor
func NewExecutor(config *config.Config, commands map[string]Command) *Executor {
	aliases := make(map[string]string)

	if config != nil {
		for _, alias := range config.Command.Aliases {
			aliases[alias.Name] = alias.Command
		}
	}

	s := &State{
		Tasks:          make([]*task.Task, 0),
		Config:         config,
		Commands:       commands,
		CommandAliases: aliases,
	}

	return &Executor{
		S: s,
	}
}

// Execute execute command string
// command string should be "cmd options"
func (e *Executor) Execute(command string) (terminate bool) {
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

	alias, ok := e.S.CommandAliases[cmd]
	if ok {
		newCommand := alias + " " + option
		fmt.Fprintf(e.S.Config.Writer, "alias %s = %s\n", cmd, alias)
		fmt.Fprintf(e.S.Config.Writer, "command: %s\n", newCommand)

		return e.Execute(newCommand)
	}

	value, ok := e.S.Commands[cmd]
	if e.S.Config != nil && e.S.Config.Writer != nil {
		if ok {
			e.S.Config.Writer.Write([]byte(cmd + " hit\n"))
		} else {
			e.S.Config.Writer.Write([]byte(cmd + " not hit\n"))
		}
	}

	if ok {
		return value.Execute(option, e.S)
	}

	// no command
	return false
}
