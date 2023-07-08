package kafka

import (
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
}
