package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
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
	automaton.Tasks, automaton.MaxTaskId = ReadTasks(automaton.Config.Paths.Task)
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
	fo, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
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
			automaton.Config.Writer.Write([]byte("append tasks to " + p + "\n"))
		}
	}

	orQuery := make([]Query, 0)
	orQuery = append(orQuery, NewNoKeyQuery("complete", make([]Query, 0), make([]Query, 0)))
	query := NewSameDayQuery("complete", time.Now(), make([]Query, 0), orQuery)
	t.writeFile(automaton.Config.Paths.Task, Ls(automaton.Tasks, query)) // write today's complete or no complete task

	automaton.Tasks, automaton.MaxTaskId = ReadTasks(automaton.Config.Paths.Task)
	return false
}

func NewSaveCommand() *SaveCommand {
	return &SaveCommand{}
}

type CompleteCommand struct {
}

func (t *CompleteCommand) completeAllSubTask(dateString string, task *Task) (completeNum int) {
	n := 0

	_, ok := task.Attributes["complete"]
	if !ok {
		// if not completed, set complete date
		task.Attributes["complete"] = dateString
		n += 1
	}

	for _, subTask := range task.SubTasks {
		n += t.completeAllSubTask(dateString, subTask)
	}
	return n
}

func (t *CompleteCommand) completeTask(taskId int, tasks []*Task) (completeTask *Task, completeNum int) {
	for _, task := range tasks {
		if task.Id == taskId {
			return task, t.completeAllSubTask(time.Now().Format(dateTimeFormat), task)
		}
		t, n := t.completeTask(taskId, task.SubTasks)
		if t != nil {
			return t, n
		}
	}
	return nil, 0
}

func (t *CompleteCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	taskId, err := strconv.Atoi(option)
	if err != nil {
		automaton.Config.Writer.Write([]byte(err.Error()))
		return false
	}

	task, n := t.completeTask(taskId, automaton.Tasks)
	if task == nil {
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("There is no Task which have task id: %d\n", taskId)))
		return false
	}

	automaton.Config.Writer.Write([]byte(fmt.Sprintf("Complete %s and %d sub tasks\n", task.Name, n)))
	return false
}

func NewCompleteCommand() *CompleteCommand {
	return &CompleteCommand{}
}

type AddTaskCommand struct {
}

func (t *AddTaskCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	nowTask, err := NewTask(option, automaton.MaxTaskId+1)
	if err != nil {
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("Create task error: %s\n", err)))
		return false
	}

	automaton.Tasks = append(automaton.Tasks, nowTask)
	automaton.MaxTaskId = nowTask.Id
	automaton.Config.Writer.Write([]byte(fmt.Sprintf("Create task: %s\n", nowTask.String(true))))
	return false
}

func NewAddTaskCommand() *AddTaskCommand {
	return &AddTaskCommand{}
}

var subTaskRegexp, _ = regexp.Compile("^([0-9])+ (.+)$")

type AddSubTaskCommand struct {
}

func (t *AddSubTaskCommand) addSubTask(taskId int, addTask *Task, tasks []*Task) (parent *Task, success bool) {
	for _, task := range tasks {
		if task.Id == taskId {
			task.SubTasks = append(task.SubTasks, addTask)
			return task, true
		}
		parent, success = t.addSubTask(taskId, addTask, task.SubTasks)
		if success {
			return parent, success
		}
	}
	return nil, false
}

func (t *AddSubTaskCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	match := subTaskRegexp.FindSubmatch([]byte(option))
	if len(match) < 3 {
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("Create Subtask error: invalid format '%s'\n", option)))
		return false
	}

	nowTask, err := NewTask(string(match[2]), automaton.MaxTaskId+1)
	if err != nil {
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("Create task error: %s\n", err)))
		return false
	}

	parentTaskId, _ := strconv.Atoi(string(match[1]))
	parent, success := t.addSubTask(parentTaskId, nowTask, automaton.Tasks)
	if success {
		automaton.MaxTaskId = nowTask.Id
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("Create SubTask:\nParent: %s\nSubTask: %s\n", parent.String(true), nowTask.String(true))))
	} else {
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("Create SubTask error: thee is no task which have :id %d\n", parentTaskId)))
	}

	return false
}

func NewAddSubTaskCommand() *AddSubTaskCommand {
	return &AddSubTaskCommand{}
}
