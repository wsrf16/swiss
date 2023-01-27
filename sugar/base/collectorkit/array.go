package collectorkit

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
