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
		log.Fatal(err)
	}
	err = InitConfig()
	if err != nil {
		log.Fatal(err)
	}
	InitLogger()

	InitWsUpgrader()
	InitRouter()
}
