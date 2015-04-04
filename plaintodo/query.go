package main

import (
	"time"
)

var dateTimeFormat = "2006-01-02 15:04"
var dateFormat = "2006-01-02"

type Query interface {
	Check(task *Task) bool
}

type QueryBase struct {
	and []Query
	or  []Query
}

func (query QueryBase) checkSubQuery(task *Task, isShow bool) bool {
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

type KeyValueQuery struct {
	*QueryBase

	key   string
	value string
}

func (query KeyValueQuery) Check(task *Task) bool {
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

type ExpireDateQuery struct {
	*QueryBase

	key   string
	value time.Time
}

func NewExpireDateQuery(key string, value time.Time, and []Query, or []Query) *ExpireDateQuery {
	return &ExpireDateQuery{
		QueryBase: &QueryBase{
			and: and,
			or:  or,
		},
		key:   key,
		value: value,
	}
}

func (query ExpireDateQuery) Check(task *Task) bool {
	if task == nil {
		return false
	}

	dateString := task.Attributes[query.key]

	var t time.Time
	t, err := time.Parse(dateTimeFormat, dateString)
	if err != nil {
		t, err = time.Parse(dateFormat, dateString)
		if err != nil {
			// no date value
			return query.checkSubQuery(task, false)
		}
	}

	return query.checkSubQuery(task, t.Before(query.value))
}
