package command

import (
	"fmt"

	"../task"
)

// AddTask create new task from input
type AddTask struct {
}

// Execute create new task from input
func (t *AddTask) Execute(option string, s *State) (terminate bool) {
	nowTask, err := task.NewTask(option, s.MaxTaskID+1)
	if err != nil {
		s.Config.Writer.Write([]byte(fmt.Sprintf("Create task error: %s\n", err)))
		return false
	}

	s.Tasks = append(s.Tasks, nowTask)
	s.MaxTaskID = nowTask.ID
	s.Config.Writer.Write([]byte(fmt.Sprintf("Create task: %s\n", nowTask.String(true))))
	return false
}

// NewAddTask return AddTask
func NewAddTask() *AddTask {
	return &AddTask{}
}
