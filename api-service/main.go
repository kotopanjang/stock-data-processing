package main

import (
	"fmt"
	"log"
	"net/http"
	"stock-data-processing/api-service/entities"
	"stock-data-processing/api-service/modules/user/handler"
	"stock-data-processing/api-service/modules/user/repository"
	"stock-data-processing/api-service/modules/user/usecase"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

func main() {
	port := viper.GetString("PORT")
	pgHost := viper.GetString("POSTGRES_HOST")         // =localhost
	pgPort := viper.GetString("POSTGRES_PORT")         // =5432
	pgUser := viper.GetString("POSTGRES_USER")         // =postgres
	pgPass := viper.GetString("POSTGRES_PASS")         // =secret
	pgDatabase := viper.GetString("POSTGRES_DATABASE") // =postgres

	// initiate go orm
	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai", pgHost, pgUser, pgPass, pgDatabase, pgPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&entities.User{})

	// initiate mux
	r := mux.NewRouter()

	// initiate repository
	userRepo := repository.NewRepository(db)

	// initiate usecase
	userUsecase := usecase.NewUsecase(userRepo)

	// initiate handler
	handler.NewHTTPHandler(r, userUsecase)

	// run the server
	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	fmt.Println("Application running on port : ", port)
	log.Fatal(srv.ListenAndServe())
}
