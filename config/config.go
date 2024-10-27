package config

import (
	"log"

	"github.com/joho/godotenv"
)

// load enviroment variables
func LoadEnv() {
	err := godotenv.Load(".env")
	//error handling
	if err != nil {
		log.Fatal("failed to load env")
	}
}
