package query

import (
	"../task"
)

// MaxLevel search task which Level smaller than MaxLevel.level
type MaxLevel struct {
	*Base
	level int
}

// NewMaxLevel return MaxLevel
func NewMaxLevel(level int, and []Query, or []Query) *MaxLevel {
	return &MaxLevel{
		Base: &Base{
			And: and,
			Or:  or,
		},
		level: level,
	}
}

// Check is search method
func (q *MaxLevel) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	return q.CheckSubQuery(task, task.Level < q.level)
}
