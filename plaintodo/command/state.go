package command

import (
	"../config"
	"../task"
)

// State is executor state data
type State struct {
	Tasks          []*task.Task
	MaxTaskID      int
	Config         *config.Config
	Commands       map[string]Command
	CommandAliases map[string]string
}
