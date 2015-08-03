package query

import (
	"../task"
)

// KeyValue search specific key and value query
type KeyValue struct {
	*Base

	key   string
	value string
}

// Check method check query have specific key and value
func (query *KeyValue) Check(task *task.Task) bool {
	if task == nil {
		return false
	}

	return query.CheckSubQuery(task, task.Attributes[query.key] == query.value)
}

// NewKeyValue is create method
func NewKeyValue(key string, value string, and []Query, or []Query) *KeyValue {
	return &KeyValue{
		Base: &Base{
			And: and,
			Or:  or,
		},
		key:   key,
		value: value,
	}
}
