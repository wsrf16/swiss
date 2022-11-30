package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

func NewSyncConsumer(addresses []string, user string, password string) (sarama.Consumer, error) {
	// consume messages after the consumer launch
	var config *sarama.Config = newKafkaConfig(user, password)
	return sarama.NewConsumer(addresses, config)
}

func Receive(topic string, consumer sarama.Consumer, do func(msg *sarama.ConsumerMessage)) error {
	defer consumer.Close()

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			log.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
			continue
		}

		func(pc sarama.PartitionConsumer) {
			defer pc.AsyncClose()
			for msg := range pc.Messages() {
				fmt.Printf("Partition:%d, Offset:%d, key:%s, value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
				do(msg)
			}
		}(pc)
	}

	return nil
}
