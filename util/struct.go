package util

import (
	"reflect"
)

func New[T any](t T) (nt T) {
	var vt = reflect.ValueOf(t)
	if vt.Kind() == reflect.Invalid || vt.IsZero() {
		return
	}

	var tt = vt.Type()
	if tt.Kind() == reflect.Pointer {
		return reflect.New(tt.Elem()).Interface().(T)
	} else {
		return reflect.New(tt).Elem().Interface().(T)
	}
}
