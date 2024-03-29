package kafkahook

import (
	"github.com/Shopify/sarama"
	"github.com/wsrf16/swiss/module/kafka"
	"github.com/wsrf16/swiss/sugar/base/ipkit"
	"github.com/wsrf16/swiss/sugar/logo"
)

type KafkaHook struct {
	producer  sarama.SyncProducer
	topic     string
	formatter logo.Formatter
}

func NewKafkaHook(producer sarama.SyncProducer, topic string, formatter logo.Formatter) KafkaHook {
	return KafkaHook{producer: producer, topic: topic, formatter: formatter}
}

func (h KafkaHook) Fire(entry *logo.Entry) error {
	ip, err := ipkit.GetHostIp()
	if err != nil {
		return err
	}
	entry.Data["host"] = ip.String()

	// b, err := h.formatter.Format(entry)
	// if err != nil {
	//    return err
	// }
	// byteEncoder := sarama.ByteEncoder(b)
	//
	// return kafka.Send(byteEncoder, h.topic, h.producer)
	return kafka.SendT(entry, h.topic, h.producer)
}

func (h KafkaHook) Levels() []logo.Level {
	return logo.AllLevels
}
