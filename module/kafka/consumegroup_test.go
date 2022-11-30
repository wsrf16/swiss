package kafka

import (
	"github.com/Shopify/sarama"
	"testing"
)

func TestConsumeGroup(t *testing.T) {
	group, err := NewSyncConsumerGroup([]string{"mecs.com:9092"}, "group-a", "", "")
	if err != nil {
		t.Error(err)
	}

	err = ReceiveGroup(group, []string{"application-log"}, func(msg *sarama.ConsumerMessage) {
		s := string(msg.Value)
		println(s)
	})
	if err != nil {
		t.Error(err)
	}
}
