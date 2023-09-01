package pubsub_test

import (
	"testing"
	"time"

	"github.com/Shopify/sarama"

	"stock-data-processing/pkg/pubsub"
	"stock-data-processing/pkg/pubsub/mocks"
)

// func TestNewSaramaKafkaConsumerGroupFullConfigAdapter_Success(t *testing.T) {
// 	subscriber, err := pubsub.NewSaramaKafkaConsumerGroupFullConfigAdapter(
// 		[]string{"kafka1.com", "kafka2.com", "kafka3.com"}, "test-group", []string{"test-topic"},
// 		pubsub.NewDefaultSaramaConsumerGroupHandler("service-test", &mocks.EventHandler{}, &mocks.DLQHandler{}),
// 		sarama.NewConfig(),
// 	)
// 	assert.Error(t, err)
// 	assert.Nil(t, subscriber)
// }

func TestSaramaKafkaConsumserGroupAdapter_Success(t *testing.T) {
	cgHandler := pubsub.NewDefaultSaramaConsumerGroupHandler("service-test", &mocks.EventHandler{}, &mocks.DLQHandler{})
	topics := []string{"test-topic"}

	cg := new(mocks.SaramaConsumerGroup)
	cg.On("Consume", Anything, Anything, Anything).Return(nil)
	// cg.On("Consume", Anything, Anything, Anything).Return(sarama.ErrOutOfBrokers)
	cg.On("Close").Return(nil)

	subscriber := pubsub.NewSaramaKafkaConsumserGroupAdapter(&pubsub.SaramaKafkaConsumserGroupAdapterConfig{
		ConsumerGroupClient:  cg,
		ConsumerGroupHandler: cgHandler,
		Topics:               topics,
	})

	subscriber.Subscribe()
	<-time.After(time.Millisecond * 10)
	subscriber.Close()

	cg.AssertExpectations(t)
}

func TestSaramaKafkaConsumserGroupAdapter_ConsumeError(t *testing.T) {
	cgHandler := pubsub.NewDefaultSaramaConsumerGroupHandler("service-test", &mocks.EventHandler{}, &mocks.DLQHandler{})
	topics := []string{"test-topic"}

	cg := new(mocks.SaramaConsumerGroup)
	// cg.On("Consume", Anything, Anything, Anything).Return(nil)
	cg.On("Consume", Anything, Anything, Anything).Return(sarama.ErrOutOfBrokers)
	cg.On("Close").Return(nil)

	subscriber := pubsub.NewSaramaKafkaConsumserGroupAdapter(&pubsub.SaramaKafkaConsumserGroupAdapterConfig{
		ConsumerGroupClient:  cg,
		ConsumerGroupHandler: cgHandler,
		Topics:               topics,
	})

	subscriber.Subscribe()
	<-time.After(time.Millisecond * 10)
	subscriber.Close()

	cg.AssertExpectations(t)
}

func TestSaramaKafkaConsumserGroupAdapter_ClosingError(t *testing.T) {
	cgHandler := pubsub.NewDefaultSaramaConsumerGroupHandler("service-test", &mocks.EventHandler{}, &mocks.DLQHandler{})
	topics := []string{"test-topic"}

	cg := new(mocks.SaramaConsumerGroup)
	cg.On("Consume", Anything, Anything, Anything).Return(nil)
	cg.On("Close").Return(sarama.ErrBrokerNotAvailable)

	subscriber := pubsub.NewSaramaKafkaConsumserGroupAdapter(&pubsub.SaramaKafkaConsumserGroupAdapterConfig{
		ConsumerGroupClient:  cg,
		ConsumerGroupHandler: cgHandler,
		Topics:               topics,
	})

	subscriber.Subscribe()
	<-time.After(time.Millisecond * 10)
	subscriber.Close()

	cg.AssertExpectations(t)
}
