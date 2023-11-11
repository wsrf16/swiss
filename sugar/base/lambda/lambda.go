package lambda

import (
	"time"
)

func Loop(always bool, do func()) {
	for first := true; first; first = always {
		do()
	}
}

//	func LoopReturn[T any](always bool, do func() T) T {
//		if always {
//			for {
//				do()
//			}
//		} else {
//			return do()
//		}
//	}
func LoopReturn[T any](always bool, do func() T) T {
	var t T
	for b := true; b; b = always {
		t = do()
	}
	return t
}

func BlockSelect(stop chan int) {
	//stop := make(chan int)
For:
	for {
		select {
		case <-stop:
			break For
		default:
			time.Sleep(time.Millisecond * 20)
		}
	}
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
