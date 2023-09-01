package pubsub

import (
	"context"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

// SaramaKafkaConsumserGroupAdapterConfig is a configuration.
//
// FIELDS:
//
// `ConsumerGroupClient` is client that returned from `sarama.NewConsumerGroup()`.
//
// `ConsumerGroupHandler` is an implementation of `sarama.ConsumerGroupHandler`.
//
// `Topics` is a kafka topic to be subscribed.
type SaramaKafkaConsumserGroupAdapterConfig struct {
	ConsumerGroupClient  sarama.ConsumerGroup
	ConsumerGroupHandler sarama.ConsumerGroupHandler
	Topics               []string
}

// SaramaKafkaConsumserGroupAdapter is an adapter for pubsub's subcriber
type SaramaKafkaConsumserGroupAdapter struct {
	closeChan chan struct{}
	config    *SaramaKafkaConsumserGroupAdapterConfig
}

// NewSaramaKafkaConsumerGroupFullConfigAdapter is constructor that immediately returns subscriber.

// NewSaramaKafkaConsumserGroupAdapter is a constructor
//
// This Constructor is deprecated and use `NewSaramaKafkaConsumerGroupFullConfigAdapter` instead.
func NewSaramaKafkaConsumserGroupAdapter(config *SaramaKafkaConsumserGroupAdapterConfig) Subscriber {
	closeChan := make(chan struct{}, 1)
	return &SaramaKafkaConsumserGroupAdapter{closeChan, config}
}

// Subscribe will consume the published message
func (skcga *SaramaKafkaConsumserGroupAdapter) Subscribe() {
	go func() {
	POLL:
		for {
			select {
			case <-skcga.closeChan:
				break POLL
			default:
				err := skcga.config.ConsumerGroupClient.Consume(context.Background(), skcga.config.Topics, skcga.config.ConsumerGroupHandler)
				if err != nil {
					log.Error().Msgf("[Sarama] %s", err.Error())
				}
			}
		}
	}()
}

// Close will stop the kafka consumer
func (skcga *SaramaKafkaConsumserGroupAdapter) Close() (err error) {
	defer close(skcga.closeChan)

	skcga.closeChan <- struct{}{}

	if err = skcga.config.ConsumerGroupClient.Close(); err != nil {
		log.Error().Msgf("[Sarama] Consumer is closed with error. | %s", err.Error())
		return
	}

	log.Info().Msg("[Sarama] Consumer is gracefully shut down.")
	return
}

// SaramaConsumerGroup is an interface that purposed for mock creation for unit testing.
// Do not use this for an implementation.
type SaramaConsumerGroup interface {
	Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error
	Errors() <-chan error
	Close() error
	Pause(topicPartitions map[string][]int32)
	PauseAll()
	Resume(topicPartitions map[string][]int32)
	ResumeAll()
}
