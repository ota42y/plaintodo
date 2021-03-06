package command

import (
	"fmt"

	"../task"
	"../util"
)

// SetAttribute set attribute to task
type SetAttribute struct {
}

func (c *SetAttribute) setAttribute(task *task.Task, attributes map[string]string) {
	for key, value := range attributes {
		task.Attributes[key] = value
	}
}

// Execute set all attribute to task which have :id
// Even if already exist, it will be overwrite.
func (c *SetAttribute) Execute(option string, s *State) (terminate bool) {
	optionMap := task.ParseOptions(" " + option)

	id, err := util.GetIntAttribute("id", optionMap)
	if err != nil {
		s.Config.Writer.Write([]byte(err.Error()))
		return false
	}
	delete(optionMap, "id")

	_, task := task.GetTask(id, s.Tasks)
	if task == nil {
		fmt.Fprintf(s.Config.Writer, "there is no exist :id %d task\n", id)
		return false
	}

	_, ok := task.Attributes["lock"]
	if ok {
		fmt.Fprintf(s.Config.Writer, "Task :id %d is locked\n", id)
		return false
	}

	c.setAttribute(task, optionMap)
	s.Config.Writer.Write([]byte(fmt.Sprintln("set attribute", task.String(true))))
	return false
}

// NewSetAttribute return SetAttribute
func NewSetAttribute() *SetAttribute {
	return &SetAttribute{}
}
