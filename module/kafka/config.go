package kafka

import (
	"crypto/sha256"
	"github.com/Shopify/sarama"
	"github.com/wsrf16/swiss/auth/scrams"
	"time"
)

func newKafkaConfig(user string, password string) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follower都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 写到随机分区中，我们默认设置32个分区
	config.Producer.Compression = sarama.CompressionSnappy    // Compress messages
	config.Producer.Flush.Frequency = 500 * time.Millisecond  // Flush batches every 500ms
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回
	config.Consumer.Group.ResetInvalidOffsets = true

	if len(user) > 0 {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = user
		config.Net.SASL.Password = password
		config.Net.SASL.Mechanism = sarama.SASLMechanism(sarama.SASLTypeSCRAMSHA256)
		config.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient { return &scrams.XDGSCRAMClient{HashGeneratorFcn: sha256.New} }
	}

	return config
}
