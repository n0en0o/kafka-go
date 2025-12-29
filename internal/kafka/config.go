package kafka

import "github.com/IBM/sarama"

const Broker = "localhost:9092"

const (
    TopicOrders        = "orders"
    TopicNotifications = "notifications"
)

func NewConfig() *sarama.Config {
    config := sarama.NewConfig()
    config.Version = sarama.V2_8_1_0
    config.Producer.Return.Successes = true
    config.Consumer.Offsets.Initial = sarama.OffsetNewest
    return config
}