package utils

import (
	"reflect"
)

func Bool[T any](value T) bool {
	switch m := any(value).(type) {
	case interface{ Bool() bool }:
		return m.Bool()
	case interface{ IsZero() bool }:
		return !m.IsZero()
	}
	return reflectValue(&value)
}

func reflectValue(vp any) bool {
	switch rv := reflect.ValueOf(vp).Elem(); rv.Kind() {
	case reflect.Map, reflect.Slice:
		return rv.Len() != 0
	default:
		is := rv.IsZero()
		return !is
	}
}

// OR 取第一个非空值
func OR[T any](conds ...T) T {
	for _, cond := range conds {
		if Bool(cond) {
			return cond
		}
	}
	return conds[len(conds)-1]
}

// IfElse 条件判断
func IfElse[T, U any](flag T, ifValue U, elseValue U) U {
	if Bool(flag) {
		return ifValue
	}
	return elseValue
}
