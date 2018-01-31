package kafka

import (
	cluster "github.com/bsm/sarama-cluster"
	"github.com/oemdaro/mqtt-microservices/data-service/appconfig"
)

// NewConsumer return new Kafka consumer group
func NewConsumer() (*cluster.Consumer, error) {
	// init (custom) config, enable errors and notifications
	consumerConfig := cluster.NewConfig()
	consumerConfig.Consumer.Return.Errors = true
	consumerConfig.Group.Return.Notifications = true

	// init consumer
	kafkaConfig := appconfig.Config.Kafka
	consumer, err := cluster.NewConsumer(kafkaConfig.Peers, kafkaConfig.GroupID, kafkaConfig.Topics, consumerConfig)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}
