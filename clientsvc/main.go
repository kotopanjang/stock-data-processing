package main

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"

	"stock-data-processing/clientsvc/middleware"
	"stock-data-processing/clientsvc/server"
	"stock-data-processing/clientsvc/service"
	redis_store "stock-data-processing/clientsvc/storage/redis"
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

func main() {
	// inititate redis
	log.Info().Msg("initiate redis ...")
	redisURL := config.GetString("redis.url")
	redisPassword := config.GetString("redis.password")
	rds := redis_store.NewRedis(redisURL, redisPassword)

	// initiate service
	svc := service.NewService(rds)

	// set router object
	router := mux.NewRouter()
	router.HandleFunc("/", healthCheck)

	httpHandler := middleware.CORS(router)

	// http handler
	service.NewHTTPHandler(router, *svc)

	appPort := config.GetString("api.port")
	svr := server.NewServer(httpHandler, appPort)
	svr.Start()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM, os.Interrupt)
	<-sigterm

	svr.Close()
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, `ok`); err != nil {
		log.Err(err).Msg(err.Error())
	}
}
