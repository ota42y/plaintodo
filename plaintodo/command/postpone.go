package command

import (
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
	// if task is locked, not change
	_, ok := task.Attributes["lock"]
	if ok {
		return fmt.Errorf("Task :id %d is locked", task.ID)
	}

	// get start time
	startString, ok := task.Attributes["start"]
	if !ok {
		return fmt.Errorf("task :id %d haven't start attribute, so postpone not work", task.ID)
	}

	_, ok = util.ParseTime(startString)
	if !ok {
		return fmt.Errorf("%s is invalid format, so postpone not work", startString)
	}

	postponeString := optionMap["postpone"]

	// :postpone 1 hour
	postponeData := strings.Split(postponeString, " ")
	if len(postponeData) != 2 {
		return fmt.Errorf("%s is invalid format", optionMap["postpone"])
	}

	postponeTime := util.AddDuration(time.Now(), postponeData[0], postponeData[1])
	if postponeTime != time.Unix(0, 0) {
		optionMap["postpone"] = postponeTime.Format(util.DateTimeFormat)
	} else {
		replaceString, isReplaced := fixStringDate(postponeString, time.Now())
		if !isReplaced {
			// invalid format
			return fmt.Errorf("'%s' is invalid format", optionMap["postpone"])
		}
		// today 20:00
		optionMap["postpone"] = replaceString
	}

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
