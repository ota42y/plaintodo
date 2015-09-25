package main

import (
	"errors"
	"fmt"
	"github.com/skratchdot/open-golang/open"
	"io"
	"sort"

	"./command"
	"./ls"
	"./output"
	"./task"
	"./util"
)

type LsCommand struct {
	w io.Writer
}

func (t *LsCommand) Execute(option string, s *command.State) (terminate bool) {
	var omitStrings []string
	tasks, isOmit := ls.ExecuteQuery(option, s.Tasks)
	if isOmit {
		omitStrings = s.Config.Command.Omits
	}

	output.Output(t.w, tasks, true, 0, omitStrings)
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
	var omitStrings []string
	output.Output(t.w, showTasks, true, 0, omitStrings)
	return false
}

func NewLsAllCommand(w io.Writer) *LsAllCommand {
	return &LsAllCommand{
		w: w,
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
