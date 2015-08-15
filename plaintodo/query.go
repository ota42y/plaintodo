package main

import (
	"time"

	"./query"
	"./task"
	"./util"
)

type AfterDateQuery struct {
	*query.Base

	key   string
	value time.Time
}

func NewAfterDateQuery(key string, value time.Time, and []query.Query, or []query.Query) *AfterDateQuery {
	return &AfterDateQuery{
		Base: &query.Base{
			And: and,
			Or:  or,
		},
		key:   key,
		value: value,
	}
}

func (query *AfterDateQuery) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	dateString := task.Attributes[query.key]

	t, ok := util.ParseTime(dateString)
	if ok {
		return query.CheckSubQuery(task, t.After(query.value))
	} else {
		return query.CheckSubQuery(task, false)
	}
}

type SameDayQuery struct {
	*query.Base

	key   string
	value time.Time
}

func NewSameDayQuery(key string, value time.Time, and []query.Query, or []query.Query) *SameDayQuery {
	return &SameDayQuery{
		Base: &query.Base{
			And: and,
			Or:  or,
		},
		key:   key,
		value: value,
	}
}

func (query *SameDayQuery) Check(task *task.Task) bool {
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
	} else {
		return query.CheckSubQuery(task, false)
	}
}
