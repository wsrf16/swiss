package kafka

import (
	"encoding/json"
	"github.com/Shopify/sarama"
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

func Send(msg *sarama.ProducerMessage, producer sarama.SyncProducer) error {
	defer producer.Close()

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Produce: Partion = %d, offset = %d", partition, offset)
	return nil
}

func SendB(b []byte, topic string, producer sarama.SyncProducer) error {
	msg := newProducerMessage(b, topic)
	return Send(msg, producer)
}

func SendT(t any, topic string, producer sarama.SyncProducer) error {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return SendB(b, topic, producer)
}

//func SendTT(t any, topic string, producer sarama.SyncProducer) error {
//    defer producer.Close()
//
//    s, err := jsonkit.Marshal(t)
//    if err != nil {
//        return err
//    }
//
//    return SendT([]byte(s), topic, producer)
//}
