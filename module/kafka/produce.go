package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/wsrf16/swiss/sugar/encoding/jsonkit"
	"log"
)

func NewSyncProducer(addresses []string, user string, password string) (sarama.SyncProducer, error) {
	var config *sarama.Config = newKafkaConfig(user, password)
	return sarama.NewSyncProducer(addresses, config)
}

func newProducerMessage(b []byte, topic string) *sarama.ProducerMessage {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(b),
	}
	return msg
}

func Send(b []byte, topic string, producer sarama.SyncProducer) error {
	defer producer.Close()

	msg := newProducerMessage(b, topic)
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Produce: Partion = %d, offset = %d\n", partition, offset)
	return nil
}

func SendT(t any, topic string, producer sarama.SyncProducer) error {
	defer producer.Close()

	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return SendT(b, topic, producer)
}

func SendTT(t any, topic string, producer sarama.SyncProducer) error {
	defer producer.Close()

	s, err := jsonkit.Marshal(t)
	if err != nil {
		return err
	}

	return SendT([]byte(s), topic, producer)
}
