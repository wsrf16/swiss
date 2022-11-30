package kafkahook

import (
	"github.com/wsrf16/swiss/module/kafka"
	"github.com/wsrf16/swiss/sugar/logo"
	"testing"
)

func TestNewKafkaLogrusHook(t *testing.T) {
	producer, err := kafka.NewSyncProducer([]string{"mecs.com:9092"}, "", "")
	if err != nil {
		t.Error(err.Error())
	}
	hook := KafkaHook{topic: "application-log", producer: producer}
	logo.AddHook(hook)

	logo.ErrorF("summaryyyyyyyy", err, "messageeeeeeeeeeee%v", 222222)
}
