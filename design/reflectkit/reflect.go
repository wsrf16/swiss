package reflectkit

import (
	"reflect"
)

func ConcreteNewPoint[T any]() *T {
	var t T
	var rType reflect.Type = reflect.TypeOf(t)
	var value reflect.Value = reflect.New(rType)
	return value.Interface().(*T)
	//return T(value.Interface())
}

func ConcreteZeroT[T any]() T {
	var t T
	var rType reflect.Type = reflect.TypeOf(t)
	var value reflect.Value = reflect.Zero(rType)
	return value.Interface().(T)
	//return T(value.Interface())
}
