package command

import (
	"../task"
)

// Reload reload tasks from file which select by config file.
type Reload struct {
}

// Execute execute reload
// overwrite s.Tasks and s.MaxTaskId
// option string not use
func (t *Reload) Execute(option string, s *State) (terminate bool) {
	s.Tasks, s.MaxTaskID = task.ReadTasks(s.Config.Paths.Task)
	return false
}

// NewReload return Reload
func NewReload() *Reload {
	return &Reload{}
}
