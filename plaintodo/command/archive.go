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
	"fmt"
	"path"
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

// Archive Archive task to file
type Archive struct {
}

func (a *Archive) collectCompleteDay(tasks []*task.Task, times *map[string]bool) {
	for _, task := range tasks {
		completeDateString, ok := task.Attributes["complete"]
		if ok {
			t, ok := util.ParseTime(completeDateString)
			if ok {
				str := t.Format(util.DateFormat)
				(*times)[str] = true
			}
		}

		a.collectCompleteDay(task.SubTasks, times)
	}
}

func (a *Archive) getCompleteDayList(tasks []*task.Task) []time.Time {
	allTimes := make(map[string]bool)
	a.collectCompleteDay(tasks, &allTimes)

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

func (a *Archive) appendFile(filePath string, tasks []*ls.ShowTask) (terminate bool, err error) {
	fo, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return false, err
	}
	defer fo.Close()

	output.Output(fo, tasks, false)
	return true, nil
}

func (a *Archive) archiveTasks(s *State) {
	// Archive old task to task file
	times := a.getCompleteDayList(s.Tasks)
	for _, value := range times {
		fileName := value.Format(s.Config.Archive.NameFormat) + ".txt"
		p := path.Join(s.Config.Archive.Directory, fileName)

		query := query.NewSameDay("complete", value, make([]query.Query, 0), make([]query.Query, 0))
		a.appendFile(p, ls.Ls(s.Tasks, query))
		fmt.Fprint(s.Config.Writer, "append tasks to %s\n", p)
	}
}

// Execute Archive tasks to file
func (a *Archive) Execute(option string, state *State) (terminate bool) {
	a.archiveTasks(state)
	return false
}

// NewArchive return Archive
func NewArchive() *Archive {
	return &Archive{}
}
