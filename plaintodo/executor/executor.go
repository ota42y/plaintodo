package executor

import (
	"fmt"
	"strings"

	"../command"
	"../config"
	"../task"
)

var optionSplit = " "

// Executor executor command
type Executor struct {
	S *command.State
}

// NewExecutor return Executor
func NewExecutor(config *config.Config, commands map[string]command.Command) *Executor {
	aliases := make(map[string]string)

	if config != nil {
		for _, alias := range config.Command.Aliases {
			aliases[alias.Name] = alias.Command
		}
	}

	s := &command.State{
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
