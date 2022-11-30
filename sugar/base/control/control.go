package control

func LoopAlways(always bool, do func()) {
	for first := true; first; first = always {
		do()
	}
}

func LoopAlwaysReturn[T any](always bool, do func() T) T {
	if always {
		for {
			do()
		}
	} else {
		return do()
	}
}

func If[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}
	return falseVal
}

func IfFunc[T any](condition bool, trueFunc func() T, falseFunc func() T) T {
	if condition {
		return trueFunc()
	}
	return falseFunc()
}

func IfFuncPair[T1 any, T2 any](condition bool, trueFunc func() (T1, T2), falseFunc func() (T1, T2)) (T1, T2) {
	if condition {
		return trueFunc()
	}
	return falseFunc()
}
