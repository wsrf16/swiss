package jsonkit

import (
	"fmt"
	"testing"
)

type Model struct {
	Id   int
	Name string
}

func TestMarshal(t *testing.T) {
	model1 := Model{Id: 1, Name: "tom"}
	marshal, err := Marshal(model1)
	if err != nil {
		t.Error(err.Error())
	}

	model2, err := UnmarshalSToT[Model](marshal)
	if err != nil {
		t.Error(err.Error())
	}

	fmt.Printf("id: %d; name: %s\n", model2.Id, model2.Name)
}
