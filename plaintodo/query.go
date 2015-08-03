package main

import (
	"time"

	"./query"
	"./task"
)

var dateTimeFormat = "2006-01-02 15:04"
var dateFormat = "2006-01-02"

func ParseTime(dateString string) (time.Time, bool) {
	var t time.Time
	t, err := time.Parse(dateTimeFormat+"-0700", dateString+"+0900")
	if err != nil {
		t, err = time.Parse(dateFormat+"-0700", dateString+"+0900")
		if err != nil {
			// not date value
			return t, false
		}
	}

	return t, true
}

type NoKeyQuery struct {
	*query.Base

	key string
}

func (query *NoKeyQuery) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	_, ok := task.Attributes[query.key]
	return query.CheckSubQuery(task, !ok)
}

func NewNoKeyQuery(key string, and []query.Query, or []query.Query) *NoKeyQuery {
	return &NoKeyQuery{
		Base: &query.Base{
			And: and,
			Or:  or,
		},
		key: key,
	}
}

type BeforeDateQuery struct {
	*query.Base

	key   string
	value time.Time
}

func NewBeforeDateQuery(key string, value time.Time, and []query.Query, or []query.Query) *BeforeDateQuery {
	return &BeforeDateQuery{
		Base: &query.Base{
			And: and,
			Or:  or,
		},
		key:   key,
		value: value,
	}
}

func (query *BeforeDateQuery) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	dateString := task.Attributes[query.key]
	t, ok := ParseTime(dateString)

	if ok {
		return query.CheckSubQuery(task, t.Before(query.value))
	} else {
		return query.CheckSubQuery(task, false)
	}
}

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

	t, ok := ParseTime(dateString)
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

	t, ok := ParseTime(dateString)

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

type MaxLevelQuery struct {
	*query.Base
	level int
}

func NewMaxLevelQuery(level int, and []query.Query, or []query.Query) *MaxLevelQuery {
	return &MaxLevelQuery{
		Base: &query.Base{
			And: and,
			Or:  or,
		},
		level: level,
	}
}

func (query *MaxLevelQuery) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	return query.CheckSubQuery(task, task.Level < query.level)
}

type IDQuery struct {
	*query.Base
	id int
}

func NewIDQuery(id int, and []query.Query, or []query.Query) *IDQuery {
	return &IDQuery{
		Base: &query.Base{
			And: and,
			Or:  or,
		},
		id: id,
	}
}

func (query *IDQuery) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	return query.CheckSubQuery(task, task.ID == query.id)
}
