package query

import (
	"time"

	"../task"
	"../util"
)

// BeforeDate search task which have before date value in BeforeDate.key
type BeforeDate struct {
	*Base

	key   string
	value time.Time
}

// NewBeforeDate return BeforeDate
func NewBeforeDate(key string, value time.Time, and []Query, or []Query) *BeforeDate {
	return &BeforeDate{
		Base: &Base{
			And: and,
			Or:  or,
		},
		key:   key,
		value: value,
	}
}

// Check check task by this and all sub query
func (q *BeforeDate) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	dateString := task.Attributes[q.key]
	t, ok := util.ParseTime(dateString)

	if ok {
		return q.CheckSubQuery(task, t.Before(q.value))
	}
	return q.CheckSubQuery(task, false)
}
