package kafka

import (
	"context"
	"github.com/Shopify/sarama"
	"log"
)

func NewSyncConsumerGroup(addresses []string, groupID, user string, password string) (sarama.ConsumerGroup, error) {
	// consume the message that the last unconsume message
	var config *sarama.Config = newKafkaConfig(user, password)
	return sarama.NewConsumerGroup(addresses, groupID, config)
}

func GroupReceive(consumerGroup sarama.ConsumerGroup, topics []string, handle func(message *sarama.ConsumerMessage) error) error {
	defer consumerGroup.Close()
	return consumerGroup.Consume(context.Background(), topics, &ConsumerGroupHandler{Handle: handle})
}

type ConsumerGroupHandler struct {
	//sarama.ConsumerGroupHandler
	Handle func(msg *sarama.ConsumerMessage) error
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
		err := handler.Handle(msg)
		if err != nil {
			log.Printf("Partition:%d, Offset:%d, key:%s err:%s", msg.Partition, msg.Offset, string(msg.Key), err)
		} else {
			session.MarkMessage(msg, "consumed")
			log.Printf("Partition:%d, Offset:%d, key:%s", msg.Partition, msg.Offset, string(msg.Key))
		}
	}

	return nil
}
