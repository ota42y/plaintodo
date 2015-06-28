package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetIntAttribute(name string, attributes map[string]string) (int, error) {
	str, ok := attributes[name]
	if !ok {
		return -1, errors.New(fmt.Sprintf("not set :%s\n", name))
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}

	return num, nil
}

func AddDuration(base time.Time, num string, unit string) time.Time {
	n, err := strconv.Atoi(num)
	if err != nil {
		return time.Unix(0, 0)
	}
	switch {
	case unit == "minutes":
		return base.Add(time.Duration(n) * time.Minute)
	case unit == "hour":
		return base.Add(time.Duration(n) * time.Hour)
	case unit == "day":
		return base.AddDate(0, 0, n)
	case unit == "week":
		return base.AddDate(0, 0, n*7)
	case unit == "month":
		return base.AddDate(0, n, 0)
	case unit == "year":
		return base.AddDate(n, 0, 0)
	}

	return time.Unix(0, 0)
}

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
	Output(t.w, ExecuteQuery(option, automaton.Tasks), true)
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
	MaxTaskId int
}

func (t *CompleteCommand) setNewRepeat(baseTime time.Time, task *Task) {
	repeatString, ok := task.Attributes["repeat"]
	if !ok {
		return
	}

	// every 1 day
	splits := strings.Split(repeatString, " ")
	if len(splits) != 3 {
		return
	}

	if splits[0] == "every" {
		startString, ok := task.Attributes["start"]
		if !ok {
			return
		}
		baseTime, ok = ParseTime(startString)
		if !ok {
			return
		}
	}

	newTime := AddDuration(baseTime, splits[1], splits[2])
	task.Attributes["start"] = newTime.Format(dateTimeFormat)
}

func (t *CompleteCommand) completeAllSubTask(completeDate time.Time, task *Task) (repeatTask *Task, completeNum int) {
	n := 0
	newSubTasks := make([]*Task, 0)

	for _, subTask := range task.SubTasks {
		repeatSubTask, num := t.completeAllSubTask(completeDate, subTask)
		n += num
		if repeatSubTask != nil {
			newSubTasks = append(newSubTasks, repeatSubTask)
		}
	}

	_, ok := task.Attributes["repeat"]
	if len(newSubTasks) != 0 || ok {
		repeatTask = task.Copy(t.MaxTaskId+1, false)
		repeatTask.SubTasks = newSubTasks
		t.setNewRepeat(completeDate, repeatTask)
		t.MaxTaskId += 1
	}

	_, ok = task.Attributes["complete"]
	if !ok {
		// if not completed, set complete date
		task.Attributes["complete"] = completeDate.Format(dateTimeFormat)
		n += 1
	}

	return repeatTask, n
}

func (t *CompleteCommand) completeTask(taskId int, tasks []*Task) (completeTask *Task, newTasks []*Task, completeNum int) {
	for _, task := range tasks {
		if task.Id == taskId {
			repeatTask, n := t.completeAllSubTask(time.Now(), task)
			if repeatTask != nil {
				tasks = append(tasks, repeatTask)
			}

			return task, tasks, n
		}
		completeTask, newTasks, n := t.completeTask(taskId, task.SubTasks)
		if completeTask != nil {
			task.SubTasks = newTasks
			return completeTask, tasks, n
		}
	}
	return nil, tasks, 0
}

func (t *CompleteCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	t.MaxTaskId = automaton.MaxTaskId

	taskId, err := strconv.Atoi(option)
	if err != nil {
		automaton.Config.Writer.Write([]byte(err.Error()))
		return false
	}

	task, newTasks, n := t.completeTask(taskId, automaton.Tasks)
	automaton.Tasks = newTasks
	automaton.MaxTaskId = t.MaxTaskId
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

var subTaskRegexp, _ = regexp.Compile("^([0-9]+) (.+)$")

type AddSubTaskCommand struct {
}

func (t *AddSubTaskCommand) addSubTask(taskId int, addTask *Task, tasks []*Task) (parent *Task, success bool) {
	for _, task := range tasks {
		if task.Id == taskId {
			addTask.Level = task.Level + 1
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

type SetAttributeCommand struct {
}

func (c *SetAttributeCommand) setAttribute(task *Task, attributes map[string]string) {
	for key, value := range attributes {
		task.Attributes[key] = value
	}
}

func (c *SetAttributeCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	optionMap := ParseOptions(" " + option)

	id, err := GetIntAttribute("id", optionMap)
	if err != nil {
		automaton.Config.Writer.Write([]byte(err.Error()))
		return false
	}
	delete(optionMap, "id")

	task := GetTask(id, automaton.Tasks)
	if task != nil {
		c.setAttribute(task, optionMap)
		automaton.Config.Writer.Write([]byte(fmt.Sprintln("set attribute", task.String(true))))
	} else {
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("there is no exist :id %d task\n", id)))
	}
	return false
}

func NewSetAttributeCommand() *SetAttributeCommand {
	return &SetAttributeCommand{}
}

type StartCommand struct {
	*SetAttributeCommand
}

func (c *StartCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	return c.SetAttributeCommand.Execute(option+" :start "+time.Now().Format(dateTimeFormat), automaton)
}

func NewStartCommand() *StartCommand {
	return &StartCommand{
		SetAttributeCommand: &SetAttributeCommand{},
	}
}

// postpone :id 1 :postpone 5 hour
type PostponeCommand struct {
	*SetAttributeCommand
}

func (c *PostponeCommand) postpone(task *Task, optionMap map[string]string) error {
	// get start time
	startString, ok := task.Attributes["start"]
	if !ok {
		return errors.New(fmt.Sprint("task :id ", task.Id, " haven't start attribute"))
	}

	startTime, ok := ParseTime(startString)
	if !ok {
		return errors.New(fmt.Sprint(startString, " is invalid format"))
	}

	// :postpone 1 hour
	postponeData := strings.Split(optionMap["postpone"], " ")
	if len(postponeData) != 2 {
		return errors.New(fmt.Sprint(optionMap["postpone"], " is invalid format"))
	}

	postponeTime := AddDuration(startTime, postponeData[0], postponeData[1])
	optionMap["postpone"] = postponeTime.Format(dateTimeFormat)

	c.setAttribute(task, optionMap)
	return nil
}

func (c *PostponeCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	optionMap := ParseOptions(" " + option)

	id, err := GetIntAttribute("id", optionMap)
	if err != nil {
		automaton.Config.Writer.Write([]byte(err.Error()))
		return false
	}
	delete(optionMap, "id")

	task := GetTask(id, automaton.Tasks)
	if task == nil {
		automaton.Config.Writer.Write([]byte(fmt.Sprintf("there is no exist :id %d task\n", id)))
		return false
	}

	err = c.postpone(task, optionMap)
	if err != nil {
		fmt.Fprintln(automaton.Config.Writer, err)
	} else {
		automaton.Config.Writer.Write([]byte(fmt.Sprintln("set attribute", task.String(true))))
	}

	return false
}

func NewPostponeCommand() *PostponeCommand {
	return &PostponeCommand{
		SetAttributeCommand: &SetAttributeCommand{},
	}
}
