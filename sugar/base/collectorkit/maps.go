package collectorkit

func GetMapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func GetMapKeysByOrder[K comparable, V any](m map[K]V, order int) K {
	if order <= len(m)-1 {
		i := 0
		for k := range m {
			if i == order {
				return k
			}
			i++
		}
	}
	panic("not exist")
}

func GetMapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func GetMapValuesByOrder[K comparable, V any](m map[K]V, order int) V {
	if order <= len(m)-1 {
		i := 0
		for _, v := range m {
			if i == order {
				return v
			}
		}
	}
	panic("not exist")
}
