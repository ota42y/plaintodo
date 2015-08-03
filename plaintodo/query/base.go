package query

import (
	"../task"
)

// Base have 'and' and 'or' condition
type Base struct {
	And []Query
	Or  []Query
}

// CheckSubQuery check all 'And' and 'Or' query and return
func (query *Base) CheckSubQuery(task *task.Task, isShow bool) bool {
	// If this query return true, check all and query
	// (Even if or query exist, we don't need check these.
	if isShow {
		for _, q := range query.And {
			if !q.Check(task) {
				return false
			}
		}
		return true
	}

	// If this query return false, check all or query
	// Even if and query exist, we don't need check these.
	for _, q := range query.Or {
		if q.Check(task) {
			return true
		}
	}
	return false
}

// Check do nothing, so CheckSubQuery
func (query *Base) Check(task *task.Task) bool {
	// do nothing
	return query.CheckSubQuery(task, true)
}

// NewBase return NewQueryBase
func NewBase(and []Query, or []Query) *Base {
	return &Base{
		And: and,
		Or:  or,
	}
}
