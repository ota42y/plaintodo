package main

import (
	"time"
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

type Query interface {
	Check(task *Task) bool
}

type QueryBase struct {
	and []Query
	or  []Query
}

func (query *QueryBase) checkSubQuery(task *Task, isShow bool) bool {
	// If this query return true, check all and query
	// (Even if or query exist, we don't need check these.
	if isShow {
		for _, q := range query.and {
			if !q.Check(task) {
				return false
			}
		}
		return true
	} else {
		// If this query return false, check all or query
		// Even if and query exist, we don't need check these.
		for _, q := range query.or {
			if q.Check(task) {
				return true
			}
		}
		return false
	}
}

func (query *QueryBase) Check(task *Task) bool {
	// do nothing
	return query.checkSubQuery(task, true)
}

func NewQueryBase(and []Query, or []Query) *QueryBase {
	return &QueryBase{
		and: and,
		or:  or,
	}
}

type KeyValueQuery struct {
	*QueryBase

	key   string
	value string
}

func (query *KeyValueQuery) Check(task *Task) bool {
	if task == nil {
		return false
	}

	return query.checkSubQuery(task, task.Attributes[query.key] == query.value)
}

func NewKeyValueQuery(key string, value string, and []Query, or []Query) *KeyValueQuery {
	return &KeyValueQuery{
		QueryBase: &QueryBase{
			and: and,
			or:  or,
		},
		key:   key,
		value: value,
	}
}

type NoKeyQuery struct {
	*QueryBase

	key string
}

func (query *NoKeyQuery) Check(task *Task) bool {
	if task == nil {
		return false
	}

	_, ok := task.Attributes[query.key]
	return query.checkSubQuery(task, !ok)
}

func NewNoKeyQuery(key string, and []Query, or []Query) *NoKeyQuery {
	return &NoKeyQuery{
		QueryBase: &QueryBase{
			and: and,
			or:  or,
		},
		key: key,
	}
}

type BeforeDateQuery struct {
	*QueryBase

	key   string
	value time.Time
}

func NewBeforeDateQuery(key string, value time.Time, and []Query, or []Query) *BeforeDateQuery {
	return &BeforeDateQuery{
		QueryBase: &QueryBase{
			and: and,
			or:  or,
		},
		key:   key,
		value: value,
	}
}

func (query *BeforeDateQuery) Check(task *Task) bool {
	if task == nil {
		return false
	}

	dateString := task.Attributes[query.key]
	t, ok := ParseTime(dateString)

	if ok {
		return query.checkSubQuery(task, t.Before(query.value))
	} else {
		return query.checkSubQuery(task, false)
	}
}

type AfterDateQuery struct {
	*QueryBase

	key   string
	value time.Time
}

func NewAfterDateQuery(key string, value time.Time, and []Query, or []Query) *AfterDateQuery {
	return &AfterDateQuery{
		QueryBase: &QueryBase{
			and: and,
			or:  or,
		},
		key:   key,
		value: value,
	}
}

func (query *AfterDateQuery) Check(task *Task) bool {
	if task == nil {
		return false
	}

	dateString := task.Attributes[query.key]

	t, ok := ParseTime(dateString)
	if ok {
		return query.checkSubQuery(task, t.After(query.value))
	} else {
		return query.checkSubQuery(task, false)
	}
}

type SameDayQuery struct {
	*QueryBase

	key   string
	value time.Time
}

func NewSameDayQuery(key string, value time.Time, and []Query, or []Query) *SameDayQuery {
	return &SameDayQuery{
		QueryBase: &QueryBase{
			and: and,
			or:  or,
		},
		key:   key,
		value: value,
	}
}

func (query *SameDayQuery) Check(task *Task) bool {
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

		return query.checkSubQuery(task, ok)
	} else {
		return query.checkSubQuery(task, false)
	}
}

type MaxLevelQuery struct {
	*QueryBase
	level int
}

func NewMaxLevelQuery(level int, and []Query, or []Query) *MaxLevelQuery {
	return &MaxLevelQuery{
		QueryBase: &QueryBase{
			and: and,
			or:  or,
		},
		level: level,
	}
}

func (query *MaxLevelQuery) Check(task *Task) bool {
	if task == nil {
		return false
	}

	return query.checkSubQuery(task, task.Level < query.level)
}
