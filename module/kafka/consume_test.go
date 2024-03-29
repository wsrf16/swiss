package kafka

import (
	"github.com/Shopify/sarama"
	"testing"
)

func TestReceives(t *testing.T) {
	consumer, err := NewSyncConsumer([]string{"mecs.com:9092"}, "", "")
	if err != nil {
		t.Error(err)
	}

	Receive(consumer, "application-log", func(msg *sarama.ConsumerMessage) {
		va := string(msg.Value)
		println(va)
	})
}
