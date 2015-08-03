package query

import (
	"../task"
)

// Query is query base interface
type Query interface {
	Check(task *task.Task) bool
}

// CreateBlankQueryArray return empty Query array
func CreateBlankQueryArray() []Query {
	return make([]Query, 0)
}
