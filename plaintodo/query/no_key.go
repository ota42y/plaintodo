package query

import (
	"../task"
)

// NoKey search task which haven't key
type NoKey struct {
	*Base

	key string
}

// Check search task
func (q *NoKey) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	_, ok := task.Attributes[q.key]
	return q.CheckSubQuery(task, !ok)
}

// NewNoKey return NoKey
func NewNoKey(key string, and []Query, or []Query) *NoKey {
	return &NoKey{
		Base: &Base{
			And: and,
			Or:  or,
		},
		key: key,
	}
}
