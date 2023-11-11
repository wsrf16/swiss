package kafka

import (
	"math/rand"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	producer, err := NewSyncProducer([]string{"mecs.com:9092"}, "", "")
	if err != nil {
		t.Error(err)
	}

	err = SendB([]byte(time.Now().String()), "application-log", producer)
	if err != nil {
		t.Error(err)
	}

	obj := make(map[string]interface{})
	obj["name"] = time.Now().String()
	obj["age"] = rand.Intn(100)

	err = SendT(obj, "application-log", producer)
	if err != nil {
		t.Error(err)
	}
}
