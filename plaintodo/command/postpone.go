package command

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"../task"
	"../util"
)

// Postpone postpone task
// postpone :id 1 :postpone 5 hour
type Postpone struct {
	*SetAttribute
}

// Postpone postpone task
// It shuld be private
func (c *Postpone) Postpone(task *task.Task, optionMap map[string]string) error {
	// get start time
	startString, ok := task.Attributes["start"]
	if !ok {
		return errors.New(fmt.Sprint("task :id ", task.ID, " haven't start attribute, so postpone not work"))
	}

	_, ok = util.ParseTime(startString)
	if !ok {
		return errors.New(fmt.Sprint(startString, " is invalid format, so postpone not work"))
	}

	// :postpone 1 hour
	postponeData := strings.Split(optionMap["postpone"], " ")
	if len(postponeData) != 2 {
		return errors.New(fmt.Sprint(optionMap["postpone"], " is invalid format"))
	}

	postponeTime := util.AddDuration(time.Now(), postponeData[0], postponeData[1])
	optionMap["postpone"] = postponeTime.Format(util.DateTimeFormat)

	c.setAttribute(task, optionMap)
	return nil
}

// Execute execute postpone
// Select task from option string
// Set postpone attribute
func (c *Postpone) Execute(option string, s *State) (terminate bool) {
	optionMap := task.ParseOptions(" " + option)

	id, err := util.GetIntAttribute("id", optionMap)
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

	err = c.Postpone(task, optionMap)
	if err != nil {
		fmt.Fprintln(s.Config.Writer, err)
	} else {
		s.Config.Writer.Write([]byte(fmt.Sprintln("set attribute", task.String(true))))
	}

	return false
}

// NewPostpone return Postpone
func NewPostpone() *Postpone {
	return &Postpone{
		SetAttribute: NewSetAttribute(),
	}
}
