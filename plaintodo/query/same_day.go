package query

import (
	"time"

	"../task"
	"../util"
)

// SameDay search task have same date value
type SameDay struct {
	*Base

	key   string
	value time.Time
}

// NewSameDay return SameDay
func NewSameDay(key string, value time.Time, and []Query, or []Query) *SameDay {
	return &SameDay{
		Base: &Base{
			And: and,
			Or:  or,
		},
		key:   key,
		value: value,
	}
}

// Check is checker method
func (query *SameDay) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	dateString := task.Attributes[query.key]

	t, ok := util.ParseTime(dateString)

	if ok {
		// not check time
		year, month, day := t.Date()
		y, m, d := query.value.Date()
		ok := (y == year) && (m == month) && (d == day)

		return query.CheckSubQuery(task, ok)
	}

	return query.CheckSubQuery(task, false)
}
