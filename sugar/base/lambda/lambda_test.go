package lambda

import (
	"testing"
)

func TestIf(t *testing.T) {

	a := 1
	b := 2
	t2 := If(a > b, "是", "否")
	t.Log(t2)

	IfFunc(a < b, func() string {
		return "是"
	}, func() string {
		return "否"
	})
}
