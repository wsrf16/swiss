package collectorkit

func ToPointerArray[T any](ts []T) []*T {
	ps := make([]*T, len(ts))
	for i, v := range ts {
		ps[i] = &v
	}
	return ps
}

func ToValueArray[T any](ts []*T) []T {
	ps := make([]T, len(ts))
	for i, v := range ts {
		ps[i] = *v
	}
	return ps
}
