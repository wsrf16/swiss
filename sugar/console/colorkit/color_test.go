package colorkit

import (
	"fmt"
	"testing"
)

func TestSpellColorString(t *testing.T) {
	s := Spell("There is some words.", GreenBg, Red)
	fmt.Println(s)

}
