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
	"time"

	"./command"
	"./ls"
	"./query"
	"./task"
	"./util"
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

func (t *ExitCommand) Execute(option string, s *command.State) (terminate bool) {
	return true
}

func NewExitCommand() *ExitCommand {
	return &ExitCommand{}
}

type LsCommand struct {
	w io.Writer
}

func (t *LsCommand) Execute(option string, s *command.State) (terminate bool) {
	Output(t.w, ls.ExecuteQuery(option, s.Tasks), true)
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
	showTasks := ls.Ls(s.Tasks, nil)
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
			t, ok := util.ParseTime(completeDateString)
			if ok {
				str := t.Format(util.DateFormat)
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
		t, ok := util.ParseTime(key)
		if ok {
			times = append(times, t)
		}
	}

	sort.Sort(times)
	return times
}

func (t *SaveCommand) appendFile(filePath string, tasks []*ls.ShowTask) (terminate bool, err error) {
	fo, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		return false, err
	}
	defer fo.Close()

	Output(fo, tasks, false)
	return true, nil
}

func (t *SaveCommand) writeFile(filePath string, tasks []*ls.ShowTask) (terminate bool, err error) {
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
			t.appendFile(p, ls.Ls(tasks, query))
			w.Write([]byte("append tasks to " + p + "\n"))
		}
	}
}

func (t *SaveCommand) saveToFile(tasks []*task.Task, saveFolder string) {
	orQuery := make([]query.Query, 0)
	orQuery = append(orQuery, query.NewNoKey("complete", make([]query.Query, 0), make([]query.Query, 0)))
	query := NewSameDayQuery("complete", time.Now(), make([]query.Query, 0), orQuery)
	t.writeFile(saveFolder, ls.Ls(tasks, query)) // write today's complete or no complete task
}

func (t *SaveCommand) Execute(option string, s *command.State) (terminate bool) {
	today, ok := util.ParseTime(time.Now().Format(util.DateFormat))
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

type StartCommand struct {
	*command.SetAttribute
}

func (c *StartCommand) Execute(option string, s *command.State) (terminate bool) {
	return c.SetAttribute.Execute(option+" :start "+time.Now().Format(util.DateTimeFormat), s)
}

func NewStartCommand() *StartCommand {
	return &StartCommand{
		SetAttribute: command.NewSetAttribute(),
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

	taskID, err := util.GetIntAttribute("from", m)
	if err != nil {
		s.Config.Writer.Write([]byte(err.Error()))
		return false
	}

	toID, err := util.GetIntAttribute("to", m)
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
	id, err := util.GetIntAttribute("id", optionMap)
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
	id, err := util.GetIntAttribute("id", optionMap)
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
