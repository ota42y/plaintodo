package main

import (
	"errors"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"io"
	"regexp"
	"sort"
	"strconv"
	"time"

	"./command"
	"./ls"
	"./output"
	"./task"
	"./util"
)

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
	output.Output(t.w, ls.ExecuteQuery(option, s.Tasks), true, 0)
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
	output.Output(t.w, showTasks, true, 0)
	return false
}

func NewLsAllCommand(w io.Writer) *LsAllCommand {
	return &LsAllCommand{
		w: w,
	}
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
