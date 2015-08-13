package query

import (
	"../task"
)

// ID search task which have specific id
type ID struct {
	*Base
	id int
}

// NewID return ID
func NewID(id int, and []Query, or []Query) *ID {
	return &ID{
		Base: &Base{
			And: and,
			Or:  or,
		},
		id: id,
	}
}

// Check check this and sub query
func (q *ID) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	return q.CheckSubQuery(task, task.ID == q.id)
}
