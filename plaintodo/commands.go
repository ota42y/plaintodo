package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

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
	showTasks := Ls(automaton.Tasks, NewExpireDateQuery("due", time.Now(), make([]Query, 0), make([]Query, 0)))
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

func (t *SaveCommand) Execute(option string, automaton *Automaton) (terminate bool) {
	fo, err := os.Create(automaton.Config.Paths.Task)
	if err != nil {
		automaton.Config.Writer.Write([]byte(err.Error()))
		return false
	}
	defer fo.Close()

	Output(fo, Ls(automaton.Tasks, nil), false) // write all task
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
