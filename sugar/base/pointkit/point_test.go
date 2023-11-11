package pointkit

import (
	"testing"
)

type Model struct {
	StrP *string
	IntP *int
}

func TestToPoint(t *testing.T) {
	m := new(Model)
	t.Log(m.StrP)
	m.StrP = ToPoint("abcdefg")
	t.Log(m.StrP)
}

func TestIsEmpty(t *testing.T) {
	var a1 *int
	empty1 := IsEmpty(a1)
	t.Log(empty1)

	var a_ = 5
	var a2 *int = &a_
	empty2 := IsEmpty(a2)
	t.Log(empty2)
}
