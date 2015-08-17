package command

import (
	"os"
	"sort"
	"time"

	"../ls"
	"../output"
	"../query"
	"../task"
	"../util"
)

type timeList []time.Time

// Len is implement for sort.Sort
func (l timeList) Len() int {
	return len(l)
}

// Less is implement for sort.Sort
func (l timeList) Less(i, j int) bool {
	return l[i].Before(l[j])
}

// Swap is implement for sort.Sort
func (l timeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type saveTask struct {
	FileName string
	Tasks    []*ls.ShowTask
}

type saveList []*saveTask

// Save save task to file
type Save struct {
}

func (s *Save) collectCompleteDay(tasks []*task.Task, times *map[string]bool) {
	for _, task := range tasks {
		completeDateString, ok := task.Attributes["complete"]
		if ok {
			t, ok := util.ParseTime(completeDateString)
			if ok {
				str := t.Format(util.DateFormat)
				(*times)[str] = true
			}
		}

		s.collectCompleteDay(task.SubTasks, times)
	}
}

func (s *Save) getCompleteDayList(tasks []*task.Task) []time.Time {
	allTimes := make(map[string]bool)
	s.collectCompleteDay(tasks, &allTimes)

	times := make(timeList, 0)
	for key := range allTimes {
		t, ok := util.ParseTime(key)
		if ok {
			times = append(times, t)
		}
	}

	sort.Sort(times)
	return times
}

func (s *Save) appendFile(filePath string, tasks []*ls.ShowTask) (terminate bool, err error) {
	fo, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return false, err
	}
	defer fo.Close()

	output.Output(fo, tasks, false)
	return true, nil
}

func (s *Save) writeFile(filePath string, tasks []*ls.ShowTask) (terminate bool, err error) {
	fo, err := os.Create(filePath)
	if err != nil {
		return false, err
	}
	defer fo.Close()

	output.Output(fo, tasks, false)
	return true, nil
}

/*
func (t *Save) archiveTasks(tasks []*task.Task, today time.Time, saveFolder string, filenameFormat string, w io.Writer) {
    // save old task to task file
    times := t.getCompleteDayList(tasks)
    for _, value := range times {
        if value != today {
            fileName := value.Format(filenameFormat) + ".txt"
            p := path.Join(saveFolder, fileName)

            query := NewSameDayQuery("complete", value, make([]query.Query, 0), make([]query.Query, 0))
            t.appendFile(p, ls.Ls(tasks, query))
            w.Write([]byte("append tasks to " + p + "\n"))
        }
    }
}
*/

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
	var orQuery []query.Query
	orQuery = append(orQuery, query.NewNoKey("complete", make([]query.Query, 0), make([]query.Query, 0)))
	q := query.NewSameDay("complete", time.Now(), make([]query.Query, 0), orQuery)
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
	/*
	   today, ok := util.ParseTime(time.Now().Format(util.DateFormat))
	   if !ok {
	       s.Config.Writer.Write([]byte("time format error"))
	       return false
	   }

	   t.archiveTasks(s.Tasks, today, s.Config.Archive.Directory, s.Config.Archive.NameFormat, s.Config.Writer)
	*/
	s.saveToFile(state)
	return false
}

// NewSave return Save
func NewSave() *Save {
	return &Save{}
}
