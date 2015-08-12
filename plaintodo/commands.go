package main

import (
	"errors"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"io"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"./command"
	"./query"
	"./task"
)

func GetIntAttribute(name string, attributes map[string]string) (int, error) {
	str, ok := attributes[name]
	if !ok {
		return -1, errors.New(fmt.Sprintf("not exist :%s\n", name))
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

func (t *ExitCommand) Execute(option string, s *command.State) (terminate bool) {
	return true
}

func NewExitCommand() *ExitCommand {
	return &ExitCommand{}
}

type ReloadCommand struct {
}

func (t *ReloadCommand) Execute(option string, s *command.State) (terminate bool) {
	s.Tasks, s.MaxTaskID = task.ReadTasks(s.Config.Paths.Task)
	return false
}

func NewReloadCommand() *ReloadCommand {
	return &ReloadCommand{}
}

type LsCommand struct {
	w io.Writer
}

func (t *LsCommand) Execute(option string, s *command.State) (terminate bool) {
	Output(t.w, ExecuteQuery(option, s.Tasks), true)
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

func (t *LsAllCommand) Execute(option string, s *command.State) (terminate bool) {
	showTasks := Ls(s.Tasks, nil)
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

func (t *SaveCommand) collectCompleteDay(tasks []*task.Task, times *map[string]bool) {
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

func (t *SaveCommand) getCompleteDayList(tasks []*task.Task) []time.Time {
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

func (t *SaveCommand) archiveTasks(tasks []*task.Task, today time.Time, saveFolder string, filenameFormat string, w io.Writer) {
	// save old task to task file
	times := t.getCompleteDayList(tasks)
	for _, value := range times {
		if value != today {
			fileName := value.Format(filenameFormat) + ".txt"
			p := path.Join(saveFolder, fileName)

			query := NewSameDayQuery("complete", value, make([]query.Query, 0), make([]query.Query, 0))
			t.appendFile(p, Ls(tasks, query))
			w.Write([]byte("append tasks to " + p + "\n"))
		}
	}
}

func (t *SaveCommand) saveToFile(tasks []*task.Task, saveFolder string) {
	orQuery := make([]query.Query, 0)
	orQuery = append(orQuery, NewNoKeyQuery("complete", make([]query.Query, 0), make([]query.Query, 0)))
	query := NewSameDayQuery("complete", time.Now(), make([]query.Query, 0), orQuery)
	t.writeFile(saveFolder, Ls(tasks, query)) // write today's complete or no complete task
}

func (t *SaveCommand) Execute(option string, s *command.State) (terminate bool) {
	today, ok := ParseTime(time.Now().Format(dateFormat))
	if !ok {
		s.Config.Writer.Write([]byte("time format error"))
		return false
	}

	t.archiveTasks(s.Tasks, today, s.Config.Archive.Directory, s.Config.Archive.NameFormat, s.Config.Writer)
	t.saveToFile(s.Tasks, s.Config.Paths.Task)
	return false
}

func NewSaveCommand() *SaveCommand {
	return &SaveCommand{}
}

type CompleteCommand struct {
	MaxTaskID int
}

func (t *CompleteCommand) setNewRepeat(baseTime time.Time, task *task.Task) {
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

func (c *CompleteCommand) completeAllSubTask(completeDate time.Time, t *task.Task) (repeatTask *task.Task, completeNum int) {
	n := 0
	newSubTasks := make([]*(task.Task), 0)

	for _, subTask := range t.SubTasks {
		repeatSubTask, num := c.completeAllSubTask(completeDate, subTask)
		n += num
		if repeatSubTask != nil {
			newSubTasks = append(newSubTasks, repeatSubTask)
		}
	}

	_, ok := t.Attributes["repeat"]
	if len(newSubTasks) != 0 || ok {
		repeatTask = t.Copy(c.MaxTaskID+1, false)
		delete(repeatTask.Attributes, "postpone")
		repeatTask.SubTasks = newSubTasks
		c.setNewRepeat(completeDate, repeatTask)
		c.MaxTaskID += 1
	}

	_, ok = t.Attributes["complete"]
	if !ok {
		// if not completed, set complete date
		t.Attributes["complete"] = completeDate.Format(dateTimeFormat)
		n += 1
	}

	return repeatTask, n
}

func (t *CompleteCommand) completeTask(taskID int, tasks []*task.Task) (completeTask *task.Task, newTasks []*task.Task, completeNum int) {
	for _, task := range tasks {
		if task.ID == taskID {
			repeatTask, n := t.completeAllSubTask(time.Now(), task)
			if repeatTask != nil {
				tasks = append(tasks, repeatTask)
			}

			return task, tasks, n
		}
		completeTask, newTasks, n := t.completeTask(taskID, task.SubTasks)
		if completeTask != nil {
			task.SubTasks = newTasks
			return completeTask, tasks, n
		}
	}
	return nil, tasks, 0
}

func (t *CompleteCommand) Execute(option string, s *command.State) (terminate bool) {
	t.MaxTaskID = s.MaxTaskID

	taskID, err := strconv.Atoi(option)
	if err != nil {
		s.Config.Writer.Write([]byte(err.Error()))
		return false
	}

	task, newTasks, n := t.completeTask(taskID, s.Tasks)
	s.Tasks = newTasks
	s.MaxTaskID = t.MaxTaskID
	if task == nil {
		s.Config.Writer.Write([]byte(fmt.Sprintf("There is no Task which have task id: %d\n", taskID)))
		return false
	}

	s.Config.Writer.Write([]byte(fmt.Sprintf("Complete %s and %d sub tasks\n", task.Name, n)))
	return false
}

func NewCompleteCommand() *CompleteCommand {
	return &CompleteCommand{}
}

type AddTaskCommand struct {
}

func (t *AddTaskCommand) Execute(option string, s *command.State) (terminate bool) {
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

func NewAddTaskCommand() *AddTaskCommand {
	return &AddTaskCommand{}
}

var subTaskRegexp, _ = regexp.Compile("^([0-9]+) (.+)$")

type AddSubTaskCommand struct {
}

func (t *AddSubTaskCommand) addSubTask(taskID int, addTask *task.Task, tasks []*task.Task) (parent *task.Task, success bool) {
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

func (t *AddSubTaskCommand) Execute(option string, s *command.State) (terminate bool) {
	match := subTaskRegexp.FindSubmatch([]byte(option))
	if len(match) < 3 {
		s.Config.Writer.Write([]byte(fmt.Sprintf("Create Subtask error: invalid format '%s'\n", option)))
		return false
	}

	nowTask, err := task.NewTask(string(match[2]), s.MaxTaskID+1)
	if err != nil {
		s.Config.Writer.Write([]byte(fmt.Sprintf("Create task error: %s\n", err)))
		return false
	}

	parentTaskID, _ := strconv.Atoi(string(match[1]))
	parent, success := t.addSubTask(parentTaskID, nowTask, s.Tasks)
	if success {
		s.MaxTaskID = nowTask.ID
		s.Config.Writer.Write([]byte(fmt.Sprintf("Create SubTask:\nParent: %s\nSubTask: %s\n", parent.String(true), nowTask.String(true))))
	} else {
		s.Config.Writer.Write([]byte(fmt.Sprintf("Create SubTask error: thee is no task which have :id %d\n", parentTaskID)))
	}

	return false
}

func NewAddSubTaskCommand() *AddSubTaskCommand {
	return &AddSubTaskCommand{}
}

type SetAttributeCommand struct {
}

func (c *SetAttributeCommand) setAttribute(task *task.Task, attributes map[string]string) {
	for key, value := range attributes {
		task.Attributes[key] = value
	}
}

func (c *SetAttributeCommand) Execute(option string, s *command.State) (terminate bool) {
	optionMap := task.ParseOptions(" " + option)

	id, err := GetIntAttribute("id", optionMap)
	if err != nil {
		s.Config.Writer.Write([]byte(err.Error()))
		return false
	}
	delete(optionMap, "id")

	_, task := task.GetTask(id, s.Tasks)
	if task != nil {
		c.setAttribute(task, optionMap)
		s.Config.Writer.Write([]byte(fmt.Sprintln("set attribute", task.String(true))))
	} else {
		s.Config.Writer.Write([]byte(fmt.Sprintf("there is no exist :id %d task\n", id)))
	}
	return false
}

func NewSetAttributeCommand() *SetAttributeCommand {
	return &SetAttributeCommand{}
}

type StartCommand struct {
	*SetAttributeCommand
}

func (c *StartCommand) Execute(option string, s *command.State) (terminate bool) {
	return c.SetAttributeCommand.Execute(option+" :start "+time.Now().Format(dateTimeFormat), s)
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

func (c *PostponeCommand) postpone(task *task.Task, optionMap map[string]string) error {
	// get start time
	startString, ok := task.Attributes["start"]
	if !ok {
		return errors.New(fmt.Sprint("task :id ", task.ID, " haven't start attribute, so postpone not work"))
	}

	_, ok = ParseTime(startString)
	if !ok {
		return errors.New(fmt.Sprint(startString, " is invalid format, so postpone not work"))
	}

	// :postpone 1 hour
	postponeData := strings.Split(optionMap["postpone"], " ")
	if len(postponeData) != 2 {
		return errors.New(fmt.Sprint(optionMap["postpone"], " is invalid format"))
	}

	postponeTime := AddDuration(time.Now(), postponeData[0], postponeData[1])
	optionMap["postpone"] = postponeTime.Format(dateTimeFormat)

	c.setAttribute(task, optionMap)
	return nil
}

func (c *PostponeCommand) Execute(option string, s *command.State) (terminate bool) {
	optionMap := task.ParseOptions(" " + option)

	id, err := GetIntAttribute("id", optionMap)
	if err != nil {
		s.Config.Writer.Write([]byte(err.Error()))
		return false
	}
	delete(optionMap, "id")

	_, task := task.GetTask(id, s.Tasks)
	if task == nil {
		s.Config.Writer.Write([]byte(fmt.Sprintf("there is no exist :id %d task\n", id)))
		return false
	}

	err = c.postpone(task, optionMap)
	if err != nil {
		fmt.Fprintln(s.Config.Writer, err)
	} else {
		s.Config.Writer.Write([]byte(fmt.Sprintln("set attribute", task.String(true))))
	}

	return false
}

func NewPostponeCommand() *PostponeCommand {
	return &PostponeCommand{
		SetAttributeCommand: &SetAttributeCommand{},
	}
}

// move :id 1 :to 1
type MoveCommand struct {
}

func (c *MoveCommand) updateTaskLevel(level int, t *task.Task) {
	t.Level = level
	for _, subTask := range t.SubTasks {
		c.updateTaskLevel(level+1, subTask)
	}
}

func (c *MoveCommand) Execute(option string, s *command.State) (terminate bool) {
	m := task.ParseOptions(" " + option)

	taskID, err := GetIntAttribute("from", m)
	if err != nil {
		s.Config.Writer.Write([]byte(err.Error()))
		return false
	}

	toID, err := GetIntAttribute("to", m)
	if err != nil {
		s.Config.Writer.Write([]byte(err.Error()))
		return false
	}

	from, t := task.GetTask(taskID, s.Tasks)
	if t == nil {
		s.Config.Writer.Write([]byte(fmt.Sprintf("there is no exist %d task\n", taskID)))
		return false
	}

	_, toTask := task.GetTask(toID, s.Tasks)
	if toTask == nil && toID != 0 {
		s.Config.Writer.Write([]byte(fmt.Sprintf("there is no exist %d task\n", toID)))
		return false
	}

	from.RemoveSubTask(t.ID)
	if toTask != nil {
		c.updateTaskLevel(toTask.Level+1, t)
		toTask.SubTasks = append(toTask.SubTasks, t)
		fmt.Fprintf(s.Config.Writer, "task moved to sub task\nparent: %s\n", toTask.String(true))
	} else {
		c.updateTaskLevel(0, t)
		s.Tasks = append(s.Tasks, t)
		fmt.Fprintf(s.Config.Writer, "task moved to top level task\n")
	}

	return false
}

func NewMoveCommand() *MoveCommand {
	return &MoveCommand{}
}

type OpenCommand struct {
}

func (c *OpenCommand) getUrl(task *task.Task) (string, error) {
	urlString, ok := task.Attributes["url"]
	if !ok {
		return "", errors.New(fmt.Sprintf("There is no url in task:\n%s", task.String(true)))
	}

	return urlString, nil
}

func (c *OpenCommand) Execute(option string, s *command.State) (terminate bool) {
	// There is no url in task:rss :id 8

	optionMap := task.ParseOptions(" " + option)
	id, err := GetIntAttribute("id", optionMap)
	if err != nil {
		fmt.Fprintf(s.Config.Writer, "%s", err.Error())
		return false
	}

	_, task := task.GetTask(id, s.Tasks)
	if task == nil {
		fmt.Fprintf(s.Config.Writer, "There is no such task :id %d\n", id)
		return false
	}

	url, err := c.getUrl(task)
	if err != nil {
		fmt.Fprintf(s.Config.Writer, "%s\n", err.Error())
		return false
	}

	open.Run(url)
	fmt.Fprintf(s.Config.Writer, "open: %s\n", url)

	return false
}

func NewOpenCommand() *OpenCommand {
	return &OpenCommand{}
}

var evernoteRegexp, _ = regexp.Compile("^https://www.evernote.com/shard/(.+)/nl/(.+)/(.+)/?")

type NiceCommand struct {
}

func (c *NiceCommand) fixEvernoteUrl(tasks []*task.Task) int {
	count := 0
	for _, task := range tasks {
		match := evernoteRegexp.FindSubmatch([]byte(task.Attributes["url"]))
		if len(match) == 4 {
			task.Attributes["url"] = fmt.Sprintf("evernote:///view/%s/%s/%s/%s/", match[2], match[1], match[3], match[3])
			count += 1
		}
		count += c.fixEvernoteUrl(task.SubTasks)
	}
	return count
}

func (c *NiceCommand) Execute(option string, s *command.State) (terminate bool) {
	var tasks []*task.Task

	optionMap := task.ParseOptions(" " + option)
	id, err := GetIntAttribute("id", optionMap)
	if err != nil {
		// do all tasks
		tasks = s.Tasks
	} else {
		// do selected task
		_, t := task.GetTask(id, s.Tasks)
		tasks = make([]*task.Task, 1)
		tasks[0] = t
	}

	fmt.Fprintf(s.Config.Writer, "Done nice\n")

	num := c.fixEvernoteUrl(tasks)
	fmt.Fprintf(s.Config.Writer, "evernote url change %d tasks\n", num)

	return false
}

func NewNiceCommand() *NiceCommand {
	return &NiceCommand{}
}

type AliasCommand struct {
}

func (c *AliasCommand) Execute(option string, s *command.State) (terminate bool) {
	keyArray := make([]string, len(s.CommandAliases))
	i := 0
	for k := range s.CommandAliases {
		keyArray[i] = k
		i++
	}
	sort.Strings(keyArray)

	for _, key := range keyArray {
		fmt.Fprintf(s.Config.Writer, "%s = %s\n", key, s.CommandAliases[key])
	}

	return false
}

func NewAliasCommand() *AliasCommand {
	return &AliasCommand{}
}
