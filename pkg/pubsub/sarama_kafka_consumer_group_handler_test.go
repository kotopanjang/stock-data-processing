package pubsub_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Shopify/sarama"

	"stock-data-processing/pkg/pubsub"
	"stock-data-processing/pkg/pubsub/mocks"
)

const Anything = "mock.Anything"

func geConsumerMessageMock() <-chan *sarama.ConsumerMessage {
	messageChan := make(chan *sarama.ConsumerMessage, 1)
	defer close(messageChan)

	message := &sarama.ConsumerMessage{
		Headers: []*sarama.RecordHeader{
			{Key: []byte("test"), Value: []byte("header")},
		},
		Key:       []byte("test-key"),
		Value:     []byte("test-message"),
		Partition: int32(1),
		Offset:    int64(40),
		Topic:     "test-topic",
	}
	messageChan <- message
	return messageChan
}

func TestSaramaKafkaConsumerGroupHandler_SuccessProceedMessage_WithoutEventHandler_WithoutDLQHandler(t *testing.T) {
	cgSess := &mocks.SaramaConsumerGroupSession{}
	cgSess.On("MarkMessage", Anything, Anything)
	cgSess.On("Context").Return(context.TODO())

	cgClaim := &mocks.SaramaConsumerGroupClaim{}
	cgClaim.On("Messages").Return(geConsumerMessageMock())

	cgh := pubsub.NewDefaultSaramaConsumerGroupHandler("service-test", nil, nil)
	if err := cgh.Setup(cgSess); err != nil {
		t.Fatal(err)
	}
	if err := cgh.ConsumeClaim(cgSess, cgClaim); err != nil {
		t.Fatal(err)
	}
	if err := cgh.Cleanup(cgSess); err != nil {
		t.Fatal(err)
	}

	cgSess.AssertExpectations(t)
	cgClaim.AssertExpectations(t)
}

func TestSaramaKafkaConsumerGroupHandler_SuccessProceedMessage_WithEventHandler_WithoutDLQHandler(t *testing.T) {
	eventHandler := &mocks.EventHandler{}
	eventHandler.On("Handle", Anything, Anything).Return(nil)

	cgSess := &mocks.SaramaConsumerGroupSession{}
	cgSess.On("MarkMessage", Anything, Anything)
	cgSess.On("Context").Return(context.TODO())

	cgClaim := &mocks.SaramaConsumerGroupClaim{}
	cgClaim.On("Messages").Return(geConsumerMessageMock())

	cgh := pubsub.NewDefaultSaramaConsumerGroupHandler("service-test", eventHandler, nil)
	if err := cgh.Setup(cgSess); err != nil {
		t.Fatal(err)
	}
	if err := cgh.ConsumeClaim(cgSess, cgClaim); err != nil {
		t.Fatal(err)
	}
	if err := cgh.Cleanup(cgSess); err != nil {
		t.Fatal(err)
	}

	cgSess.AssertExpectations(t)
	cgClaim.AssertExpectations(t)
	eventHandler.AssertExpectations(t)
}

func TestSaramaKafkaConsumerGroupHandler_ErrorProceedMessage_WithEventHandler_WithoutDLQHandler(t *testing.T) {
	eventHandler := &mocks.EventHandler{}
	eventHandler.On("Handle", Anything, Anything).Return(fmt.Errorf("error"))

	cgSess := &mocks.SaramaConsumerGroupSession{}
	cgSess.On("MarkMessage", Anything, Anything)
	cgSess.On("Context").Return(context.TODO())

	cgClaim := &mocks.SaramaConsumerGroupClaim{}
	cgClaim.On("Messages").Return(geConsumerMessageMock())

	cgh := pubsub.NewDefaultSaramaConsumerGroupHandler("service-test", eventHandler, nil)
	if err := cgh.Setup(cgSess); err != nil {
		t.Fatal(err)
	}
	if err := cgh.ConsumeClaim(cgSess, cgClaim); err != nil {
		t.Fatal(err)
	}
	if err := cgh.Cleanup(cgSess); err != nil {
		t.Fatal(err)
	}

	cgSess.AssertExpectations(t)
	cgClaim.AssertExpectations(t)
	eventHandler.AssertExpectations(t)
}

func TestSaramaKafkaConsumerGroupHandler_ErrorProceedMessage_WithEventHandler_WithDLQHandler(t *testing.T) {
	eventHandler := &mocks.EventHandler{}
	eventHandler.On("Handle", Anything, Anything).Return(fmt.Errorf("error"))

	// Send(ctx context.Context, topic string, key string, headers MessageHeaders, message []byte) (err error)
	publisher := &mocks.Publisher{}
	publisher.On("Send", Anything, Anything, Anything, Anything, Anything).Return(nil)

	dlqHandler := pubsub.NewDLQHandlerAdapter("dlq-test-topic", publisher)

	cgSess := &mocks.SaramaConsumerGroupSession{}
	cgSess.On("MarkMessage", Anything, Anything)
	cgSess.On("Context").Return(context.TODO())

	cgClaim := &mocks.SaramaConsumerGroupClaim{}
	cgClaim.On("Messages").Return(geConsumerMessageMock())

	cgh := pubsub.NewDefaultSaramaConsumerGroupHandler("service-test", eventHandler, dlqHandler)
	if err := cgh.Setup(cgSess); err != nil {
		t.Fatal(err)
	}
	if err := cgh.ConsumeClaim(cgSess, cgClaim); err != nil {
		t.Fatal(err)
	}
	if err := cgh.Cleanup(cgSess); err != nil {
		t.Fatal(err)
	}

	cgSess.AssertExpectations(t)
	cgClaim.AssertExpectations(t)
	eventHandler.AssertExpectations(t)
	publisher.AssertExpectations(t)
}
