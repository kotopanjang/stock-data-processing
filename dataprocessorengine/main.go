package main

import (
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"stock-data-processing/dataprocessorengine/engine"
	redis_egn "stock-data-processing/dataprocessorengine/storage/redis"
	"stock-data-processing/pkg/pubsub"
)

const (
	engineName                 = "dataprocessorengine"
	headerPub                  = "fileprocessorengine"
	topicFileProcessing string = "raw-data-ready-to-process"
)

var (
	config *viper.Viper
)

func init() {
	config = viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}
}

func initSaramaConfig(config *viper.Viper) *sarama.Config {
	kafkaUser := config.GetString("kafka.username")
	kafkaPassword := config.GetString("kafka.password")
	kafkaSsl := config.GetBool("kafka.ssl")
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_1_0_0
	if kafkaUser != "" {
		saramaConfig.Net.SASL.User = kafkaUser
		saramaConfig.Net.SASL.Password = kafkaPassword
		saramaConfig.Net.SASL.Enable = true
	}
	saramaConfig.Net.TLS.Enable = kafkaSsl

	// consumer config
	saramaConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	// producer config
	saramaConfig.Producer.Retry.Backoff = time.Millisecond * 500

	return saramaConfig
}

func initSaramaConsumer(saramaConfig *sarama.Config, saramaAddr []string, handler pubsub.EventHandler) pubsub.Subscriber {
	// setup consumer group handler
	paidSingleFuelDirectOrderConsumerGroupHandler := pubsub.NewDefaultSaramaConsumerGroupHandler(engineName, handler, nil)

	// setup consumer group
	paidSingleFuelDirectOrderConsumerGroup, err := sarama.NewConsumerGroup(saramaAddr, engineName, saramaConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("error init sarama consumer")
	}

	// set subscriber
	subs := pubsub.NewSaramaKafkaConsumserGroupAdapter(&pubsub.SaramaKafkaConsumserGroupAdapterConfig{
		ConsumerGroupClient:  paidSingleFuelDirectOrderConsumerGroup,
		ConsumerGroupHandler: paidSingleFuelDirectOrderConsumerGroupHandler,
		Topics:               []string{topicFileProcessing},
	})

	return subs
}

func main() {
	log.Info().Msg("Starting ...")

	// initiate kafka producer
	// set publisher object
	log.Info().Msg("initiate kafka consumer ...")
	_saramaAddr := config.GetString("kafka.brokers")
	saramaAddr := strings.Split(_saramaAddr, ",")
	saramaConfig := initSaramaConfig(config)
	// publisher := initSaramaProducer(saramaConfig, saramaAddr)

	// init redis
	log.Info().Msg("initiate redis ...")
	redisURL := config.GetString("redis.url")
	redisPassword := config.GetString("redis.password")
	rds := redis_egn.NewRedis(redisURL, redisPassword)

	// init Engine
	egn, _ := engine.NewEngine(rds)

	// init sarama subscriber
	saramaSubs := initSaramaConsumer(saramaConfig, saramaAddr, egn)
	saramaSubs.Subscribe()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, os.Interrupt)
	<-sigterm

	if err := saramaSubs.Close(); err != nil {
		log.Err(err).Msg(err.Error())
		return
	}
}
