package command

import (
	"../util"
	"time"
)

// Start set nod time to task
type Start struct {
	*SetAttribute
}

// Execute set nod time to task
func (c *Start) Execute(option string, s *State) (terminate bool) {
	return c.SetAttribute.Execute(option+" :start "+time.Now().Format(util.DateTimeFormat), s)
}

// NewStart return Start
func NewStart() *Start {
	return &Start{
		SetAttribute: NewSetAttribute(),
	}
}
