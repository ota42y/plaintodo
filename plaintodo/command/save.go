package command

import (
	"os"

	"../ls"
	"../output"
	"../query"
)

type saveTask struct {
	FileName string
	Tasks    []*ls.ShowTask
}

type saveList []*saveTask

// Save save task to file
type Save struct {
}

func (s *Save) writeFile(filePath string, tasks []*ls.ShowTask) (terminate bool, err error) {
	fo, err := os.Create(filePath)
	if err != nil {
		return false, err
	}
	defer fo.Close()

	var omitStrings []string
	output.Output(fo, tasks, false, 0, omitStrings)
	return true, nil
}

func (s *Save) addTaskToSaveTask(tasks []*ls.ShowTask, nowFilename string, list saveList) saveList {
	var saveTasks saveTask
	saveTasks.FileName = nowFilename
	saveTasks.Tasks = make([]*ls.ShowTask, 0)

	for _, nowTask := range tasks {
		saveTasks.Tasks = append(saveTasks.Tasks, nowTask)
		fileName, ok := nowTask.Task.Attributes["subTaskFile"]
		if ok {
			list = s.addTaskToSaveTask(nowTask.SubTasks, fileName, list)
			nowTask.SubTasks = make([]*ls.ShowTask, 0)
		}
	}

	list = append(list, &saveTasks)
	return list
}

func (s *Save) getSaveTaskList(state *State) *saveList {
	q := query.NewNoKey("complete", make([]query.Query, 0), make([]query.Query, 0))
	showTasks := ls.Ls(state.Tasks, q)

	var list saveList
	ret := s.addTaskToSaveTask(showTasks, state.Config.Task.DefaultFilename, list)
	return &ret
}

func (s *Save) saveToFile(state *State) {
	saveTasks := s.getSaveTaskList(state)

	for _, data := range *saveTasks {
		s.writeFile(state.Config.Task.GetTaskFilepath(data.FileName), data.Tasks)
	}
}

// Execute save tasks to file
func (s *Save) Execute(option string, state *State) (terminate bool) {
	s.saveToFile(state)
	return false
}

// NewSave return Save
func NewSave() *Save {
	return &Save{}
}
