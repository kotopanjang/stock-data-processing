package pubsub

import (
	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

// SaramaKafkaProducerAdapterConfig is a configuration of sarama kafka adapter.
type SaramaKafkaProducerAdapterConfig struct {
	AsyncProducer sarama.AsyncProducer
}

// SaramaKafkaProducerAdapter is a concrete struct of sarma kafka adapter.
type SaramaKafkaProducerAdapter struct {
	config *SaramaKafkaProducerAdapterConfig
}

// NewSaramaKafkaProducerAdapter is a constructor.
func NewSaramaKafkaProducerAdapter(config *SaramaKafkaProducerAdapterConfig) Publisher {
	p := &SaramaKafkaProducerAdapter{
		config,
	}
	go p.run()

	return p
}

func (skpa *SaramaKafkaProducerAdapter) run() {
	for producerError := range skpa.config.AsyncProducer.Errors() {
		log.Error().Msgf("[Sarama] %s", producerError.Error())
	}
}

// Send will send the message to the brokers
func (skpa *SaramaKafkaProducerAdapter) Send(topic string, key string, headers MessageHeaders, message []byte) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			log.Error().Msgf("[Sarama] %v", r)
		}
	}()
	bunchOfRecordHeaders := make([]sarama.RecordHeader, 0)

	producerMessage := &sarama.ProducerMessage{
		Headers: bunchOfRecordHeaders,
		Key:     sarama.ByteEncoder(key),
		Topic:   topic,
		Value:   sarama.ByteEncoder(message),
	}

	channel := skpa.config.AsyncProducer.Input()
	channel <- producerMessage

	return
}

// Close will stop the producer
func (skpa *SaramaKafkaProducerAdapter) Close() (err error) {
	err = skpa.config.AsyncProducer.Close()
	log.Info().Msg("[Sarama] Producer is gracefully shutdown")
	return
}
