package command

import (
	"fmt"

	"../task"
)

// Reload reload tasks from file which select by config file.
type Reload struct {
}

// Execute execute reload
// overwrite s.Tasks and s.MaxTaskId
// option string not use
func (r *Reload) Execute(option string, s *State) (terminate bool) {
	var err error

	taskFilepath := s.Config.Task.GetDefaultTaskFilepath()
	s.Tasks, s.MaxTaskID, err = task.ReadTasks(taskFilepath, 0)
	fmt.Println(taskFilepath)
	if err != nil {
		fmt.Fprintf(s.Config.Writer, "%v\n", err)
	} else {
		r.readSubTaskFile(s)
	}

	return false
}

func incrementTaskLevel(tasks []*task.Task, level int) {
	for _, nowTask := range tasks {
		nowTask.Level = level
		incrementTaskLevel(nowTask.SubTasks, level+1)
	}
}

func (r *Reload) readSubTaskFile(s *State) {
	var err error

	for _, nowTask := range s.Tasks {
		filepath, ok := nowTask.Attributes["subTaskFile"]
		if ok {
			beforeNum := s.MaxTaskID
			taskFilepath := s.Config.Task.GetTaskFilepath(filepath)
			nowTask.SubTasks, s.MaxTaskID, err = task.ReadTasks(taskFilepath, s.MaxTaskID)
			incrementTaskLevel(nowTask.SubTasks, nowTask.Level+1)

			if err != nil {
				fmt.Fprintf(s.Config.Writer, "%v\n", err)
			} else {
				fmt.Fprintf(s.Config.Writer, "read %d tasks from %s\n", (s.MaxTaskID - beforeNum), taskFilepath)
			}
		}
	}
}

// NewReload return Reload
func NewReload() *Reload {
	return &Reload{}
}
