package kafka

import (
	"github.com/Shopify/sarama"
	"log"
)

func NewSyncConsumer(addresses []string, user string, password string) (sarama.Consumer, error) {
	// consume messages after the consumer launch
	var config *sarama.Config = newKafkaConfig(user, password)
	return sarama.NewConsumer(addresses, config)
}

func Receive(consumer sarama.Consumer, topic string, handle func(msg *sarama.ConsumerMessage)) error {
	defer consumer.Close()

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to start consumer for partition %d: %s", partition, err)
			continue
		}

		func(pc sarama.PartitionConsumer) {
			defer pc.AsyncClose()
			for msg := range pc.Messages() {
				log.Printf("Partition:%d, Offset:%d, key:%s", msg.Partition, msg.Offset, string(msg.Key))
				handle(msg)
			}
		}(pc)
	}

	return nil
}
