package timekit

import (
	"fmt"
	"github.com/wsrf16/swiss/sugar/encoding/jsonkit"
	"testing"
	"time"
)

type Person struct {
	Birthday *NormalTime `json:"birthday"`
	Name     string      `json:"name"`
}

func TestMarshal(t *testing.T) {
	normalTime := NormalTime(time.Now())
	var p = Person{Name: "leon", Birthday: &normalTime}
	data, err := jsonkit.Marshal(p)
	if err != nil {
		fmt.Println("marshal err", err)
		return
	}
	t.Log(data)

	// 反序列化
	src := `{"birthday":"2020-05-26 20:20:44","name":"leon"}`
	pp, _ := jsonkit.UnmarshalSToT[Person](src)
	t.Log(pp)
}
