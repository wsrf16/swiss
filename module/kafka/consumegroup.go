package kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
)

func NewSyncConsumerGroup(addresses []string, groupID, user string, password string) (sarama.ConsumerGroup, error) {
	// consume the message that the last unconsume message
	var config *sarama.Config = newKafkaConfig(user, password)
	return sarama.NewConsumerGroup(addresses, groupID, config)
}

func ReceiveGroup(consumerGroup sarama.ConsumerGroup, topics []string, handle func(message *sarama.ConsumerMessage)) error {
	defer consumerGroup.Close()
	return consumerGroup.Consume(context.Background(), topics, &ConsumerGroupHandler{Do: func(msg *sarama.ConsumerMessage) {
		handle(msg)
	}})
}

type ConsumerGroupHandler struct {
	//sarama.ConsumerGroupHandler
	Do func(msg *sarama.ConsumerMessage)
}

func (*ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (*ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (handler *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Partition:%d, Offset:%d, key:%s, value:%s\n", msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		handler.Do(msg)
		session.MarkMessage(msg, "")
	}

	return nil
}
