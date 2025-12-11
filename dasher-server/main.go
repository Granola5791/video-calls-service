package main

import (
	"log"
)

func main() {
	err := InitConfig()
	if err != nil {
		log.Fatal(err)
	}
}
