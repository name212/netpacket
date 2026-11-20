// Copyright 2026
// license that can be found in the LICENSE file.

package value

import "reflect"

func IsNil(value any) bool {
	// copy from https://github.com/name212/govalue/blob/c6f1c3ba706eed264adfcac58cfba34fbfd9fe12/is_nil.go#L13
	// to save no deps
	iv := reflect.ValueOf(value)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	//nolint:govet
	case reflect.Interface, reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Chan:
		return iv.IsNil()
	default:
		return false
	}
}
