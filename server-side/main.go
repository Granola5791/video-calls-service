package main

import (
	"log"

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
	err = InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	err = InitDatabaseConnection()
	if err != nil {
		log.Fatal(err)
	}
	InitRouter()
}