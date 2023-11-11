package collectorkit

import "math/rand"

func ToPointerSlice[T any](ts []T) []*T {
	ps := make([]*T, len(ts))
	for i, v := range ts {
		ps[i] = &v
	}
	return ps
}

func ToValueSlice[T any](ts []*T) []T {
	ps := make([]T, len(ts))
	for i, v := range ts {
		ps[i] = *v
	}
	return ps
}

func NewSlice(ts ...any) []any {
	slice := make([]any, len(ts))
	for _, t := range ts {
		slice = append(slice, t)
	}
	return slice
}

func PickRandom[T any](arrays []T) T {
	return arrays[rand.Intn(len(arrays))]
}

func ConvertItem[T any, U any](ts []T, covert func(T) U) []U {
	us := make([]U, len(ts))
	for _, t := range ts {
		u := covert(t)
		us = append(us, u)
	}
	return us
}
