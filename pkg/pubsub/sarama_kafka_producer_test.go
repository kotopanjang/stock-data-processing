package pubsub_test

import (
	"testing"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"

	"stock-data-processing/pkg/pubsub"
)

func TestSaramaKafkaProducer(t *testing.T) {
	t.Run("when the message is successfully sent and delivered", func(t *testing.T) {
		saramaProducer := mocks.NewAsyncProducer(t, nil)
		saramaProducer.ExpectInputAndSucceed()

		publisher := pubsub.NewSaramaKafkaProducerAdapter(&pubsub.SaramaKafkaProducerAdapterConfig{
			AsyncProducer: saramaProducer,
		})

		headers := pubsub.MessageHeaders{}
		headers.Add("test", "header")

		err := publisher.Send("test-topic", "test-key", headers, []byte("Hola"))
		if err != nil {
			t.Errorf("want nil; got %v", err)
		}

		publisher.Close()
	})

	t.Run("when the connection is timeout", func(t *testing.T) {
		saramaProducer := mocks.NewAsyncProducer(t, nil)
		saramaProducer.ExpectInputAndFail(sarama.ErrRequestTimedOut)

		publisher := pubsub.NewSaramaKafkaProducerAdapter(&pubsub.SaramaKafkaProducerAdapterConfig{
			AsyncProducer: saramaProducer,
		})

		headers := pubsub.MessageHeaders{}
		headers.Add("test", "header")

		err := publisher.Send("test-topic", "test-key", headers, []byte("Hola"))
		if err != nil {
			t.Errorf("want nil; got %v", err)
		}
		publisher.Close()
	})

	t.Run("when the message channel is already closed", func(t *testing.T) {
		saramaProducer := mocks.NewAsyncProducer(t, nil)

		publisher := pubsub.NewSaramaKafkaProducerAdapter(&pubsub.SaramaKafkaProducerAdapterConfig{
			AsyncProducer: saramaProducer,
		})

		headers := pubsub.MessageHeaders{}
		headers.Add("test", "header")

		publisher.Close()
		err := publisher.Send("test-topic", "test-key", headers, []byte("Hola"))
		if err != nil {
			t.Errorf("want nil; got %v", err)
		}
	})
}
