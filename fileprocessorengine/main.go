package main

import (
	"context"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"stock-data-processing/fileprocessorengine/engine"
	"stock-data-processing/fileprocessorengine/filereader"
	"stock-data-processing/pkg/pubsub"
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

func initSaramaProducer(saramaConfig *sarama.Config, saramaAddr []string) pubsub.Publisher {
	saramaAsyncProducer, err := sarama.NewAsyncProducer(
		saramaAddr,
		saramaConfig,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("error init sarama")
	}
	publisher := pubsub.NewSaramaKafkaProducerAdapter(&pubsub.SaramaKafkaProducerAdapterConfig{
		AsyncProducer: saramaAsyncProducer,
	})
	return publisher
}

func main() {
	log.Info().Msg("Starting ...")

	// initiate kafka producer
	// set publisher object
	log.Info().Msg("initiate kafka producer ...")
	_saramaAddr := config.GetString("kafka.brokers")
	saramaAddr := strings.Split(_saramaAddr, ",")
	saramaConfig := initSaramaConfig(config)
	publisher := initSaramaProducer(saramaConfig, saramaAddr)

	// initiate file reader
	log.Info().Msg("initiate file reader ...")
	fileLocation := config.GetString("file.raw")
	successLocation := config.GetString("file.success")
	failLocation := config.GetString("file.fail")
	fr, err := filereader.NewFileReader(fileLocation, successLocation, failLocation, []string{"ndjson"})
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	// initiate engine
	log.Info().Msg("initiate engine ...")
	enableMoveFile := config.GetBool("app.enableMoveFile")
	egn, err := engine.NewEngine(fr, publisher, enableMoveFile)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}

	// process file
	ctx := context.Background()
	result := egn.ProcessFile(ctx)
	log.Info().Msgf("[File] count   : %v", result.FileCount)
	log.Info().Msgf("[File] success count     : %v", result.SuccessCount)
	log.Info().Msgf("[File] failed to process : %v", len(result.FailFile))
	log.Info().Msgf("[Data] count   : %v", result.DataCount)
	log.Info().Msgf("[Data] success count     : %v", result.DataSuccessCount)
	log.Info().Msgf("[Data] failed to process : %v", len(result.FailData))

	defer func() {
		err := publisher.Close()
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}
	}()
}
