package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"sort"
	"strconv"
	"time"
)

type timeList []time.Time

func (l timeList) Len() int {
	return len(l)
}

func (l timeList) Less(i, j int) bool {
	return l[i].Before(l[j])
}

func (l timeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type ExitCommand struct {
}

func (t *ExitCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	return true
}

func NewExitCommand() *ExitCommand {
	return &ExitCommand{}
}

type ReloadCommand struct {
}

func (t *ReloadCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	automaton.Tasks = ReadTasks(automaton.Config.Paths.Task)
	return false
}

func NewReloadCommand() *ReloadCommand {
	return &ReloadCommand{}
}

type LsCommand struct {
	w io.Writer
}

func (t *LsCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	showTasks := Ls(automaton.Tasks, NewBeforeDateQuery("due", time.Now(), make([]Query, 0), make([]Query, 0)))
	Output(t.w, showTasks, true)
	return false
}

func NewLsCommand(w io.Writer) *LsCommand {
	return &LsCommand{
		w: w,
	}
}

type LsAllCommand struct {
	w io.Writer
}

func (t *LsAllCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	showTasks := Ls(automaton.Tasks, nil)
	Output(t.w, showTasks, true)
	return false
}

func NewLsAllCommand(w io.Writer) *LsAllCommand {
	return &LsAllCommand{
		w: w,
	}
}

type SaveCommand struct {
}

func (t *SaveCommand) collectCompleteDay(tasks []*Task, times *map[string]bool) {
	for _, task := range tasks {
		completeDateString, ok := task.Attributes["complete"]
		if ok {
			t, ok := ParseTime(completeDateString)
			if ok {
				str := t.Format(dateFormat)
				(*times)[str] = true
			}
		}

		t.collectCompleteDay(task.SubTasks, times)
	}
}

func (t *SaveCommand) getCompleteDayList(tasks []*Task) []time.Time {
	allTimes := make(map[string]bool)
	t.collectCompleteDay(tasks, &allTimes)

	times := make(timeList, 0)
	for key, _ := range allTimes {
		t, ok := ParseTime(key)
		if ok {
			times = append(times, t)
		}
	}

	sort.Sort(times)
	return times
}

func (t *SaveCommand) appendFile(filePath string, tasks []*ShowTask) (terminate bool, err error) {
	fo, err := os.OpenFile(filePath, os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		return false, err
	}
	defer fo.Close()

	Output(fo, tasks, false)
	return true, nil
}

func (t *SaveCommand) writeFile(filePath string, tasks []*ShowTask) (terminate bool, err error) {
	fo, err := os.Create(filePath)
	if err != nil {
		return false, err
	}
	defer fo.Close()

	Output(fo, tasks, false)
	return true, nil
}

func (t *SaveCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	today, ok := ParseTime(time.Now().Format(dateFormat))
	if !ok {
		automaton.Config.Writer.Write([]byte("time format error"))
		return false
	}

	// save old task to task file
	times := t.getCompleteDayList(automaton.Tasks)
	for _, value := range times {
		if value != today {
			fileName := value.Format(automaton.Config.Archive.NameFormat) + ".txt"
			p := path.Join(automaton.Config.Archive.Directory, fileName)

			query := NewSameDayQuery("complete", value, make([]Query, 0), make([]Query, 0))
			t.appendFile(p, Ls(automaton.Tasks, query))
		}
	}

	orQuery := make([]Query, 0)
	orQuery = append(orQuery, NewNoKeyQuery("complete", make([]Query, 0), make([]Query, 0)))
	query := NewSameDayQuery("complete", time.Now(), make([]Query, 0), orQuery)
	t.writeFile(automaton.Config.Paths.Task, Ls(automaton.Tasks, query)) // write today's complete or no complete task
	return false
}

func NewSaveCommand() *SaveCommand {
	return &SaveCommand{}
}

type CompleteCommand struct {
}

func (t *CompleteCommand) completeTask(taskId int, tasks []*Task) (isComplete bool) {
	for _, task := range tasks {
		if task.Id == taskId {
			task.Attributes["complete"] = time.Now().Format(dateTimeFormat)
			return true
		}
		if t.completeTask(taskId, task.SubTasks) {
			return true
		}
	}
	return false
}

func (t *CompleteCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	taskId, err := strconv.Atoi(option)
	if err != nil {
		automaton.Config.Writer.Write([]byte(err.Error()))
		return false
	}

	if !t.completeTask(taskId, automaton.Tasks) {
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("There is no Task which have task id: %d\n", taskId)))
		return false
	}

	return false
}

func NewCompleteCommand() *CompleteCommand {
	return &CompleteCommand{}
}
