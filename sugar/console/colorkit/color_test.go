package colorkit

import (
	"fmt"
	"testing"
)

func TestSpellColorString(t *testing.T) {
	s := SpellColorString("There is some words.", GreenBg, Red)
	fmt.Println(s)

}
