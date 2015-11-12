package command

import (
	"fmt"
	"strconv"

	"../task"
)

// AddTask create task and add as a parent task's sub task
type AddTask struct {
}

func (t *AddTask) addSubTask(taskID int, addTask *task.Task, tasks []*task.Task) (parent *task.Task, success bool) {
	for _, task := range tasks {
		if task.ID == taskID {
			addTask.Level = task.Level + 1
			task.SubTasks = append(task.SubTasks, addTask)
			return task, true
		}
		parent, success = t.addSubTask(taskID, addTask, task.SubTasks)
		if success {
			return parent, success
		}
	}
	return nil, false
}

// Execute create sub task from input
// subtask task name :id parentId :attribute attr
// parent task id must be set
func (t *AddTask) Execute(option string, s *State) (terminate bool) {
	nowTask, err := task.NewTask(option, s.MaxTaskID+1)
	if err != nil {
		s.Config.Writer.Write([]byte(fmt.Sprintf("Create task error: %s\n", err)))
		return false
	}

	// get parent id from input
	idStr, ok := nowTask.Attributes["id"]
	if !ok {
		// new task in top level
		s.Tasks = append(s.Tasks, nowTask)
		s.MaxTaskID = nowTask.ID
		s.Config.Writer.Write([]byte(fmt.Sprintf("Create task: %s\n", nowTask.String(true))))
		return false
	}
	// subtask

	// delete parent id from attribute
	delete(nowTask.Attributes, "id")

	parentTaskID, err := strconv.Atoi(string(idStr))
	if err != nil {
		fmt.Fprintf(s.Config.Writer, "Parent id format error %s", err.Error())
		return false
	}

	parent, success := t.addSubTask(parentTaskID, nowTask, s.Tasks)
	if success {
		s.MaxTaskID = nowTask.ID
		fmt.Fprintf(s.Config.Writer, "Create SubTask:\nParent: %s\nSubTask: %s\n", parent.String(true), nowTask.String(true))
	} else {
		fmt.Fprintf(s.Config.Writer, "Create SubTask error: thee is no task which have :id %d\n", parentTaskID)
	}

	return false
}

// NewAddTask return AddSubTask
func NewAddTask() *AddTask {
	return &AddTask{}
}
