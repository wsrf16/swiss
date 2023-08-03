package lambda

func LoopAlways(always bool, do func()) {
	for first := true; first; first = always {
		do()
	}
}

//	func LoopAlwaysReturn[T any](always bool, do func() T) T {
//		if always {
//			for {
//				do()
//			}
//		} else {
//			return do()
//		}
//	}
func LoopAlwaysReturn[T any](always bool, do func() T) T {
	var t T
	for b := true; b; b = always {
		t = do()
	}
	return t
}

func If[T any](condition bool, match, notMatch T) T {
	if condition {
		return match
	}
	return notMatch
}

func IfFunc[T any](condition bool, matchFunc func() T, notMatchFunc func() T) T {
	if condition {
		return matchFunc()
	}
	return notMatchFunc()
}

func IfFuncPair[T1 any, T2 any](condition bool, matchFunc func() (T1, T2), notMatchFunc func() (T1, T2)) (T1, T2) {
	if condition {
		return matchFunc()
	}
	return notMatchFunc()
}
