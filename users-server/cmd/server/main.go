package main

import (
	"log"

	"github.com/Granola5791/video-calls-service/internal/api"
	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/Granola5791/video-calls-service/internal/logger"
	"github.com/Granola5791/video-calls-service/internal/mywebsocket"
	"github.com/joho/godotenv"
)

func InitEnv() error {
	return godotenv.Load()
}

func main() {
	err := InitEnv()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	logger.InitLogger()
	err = db.InitDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	mywebsocket.InitWsUpgrader()
	api.InitRouter()
}
