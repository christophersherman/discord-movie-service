package config

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// Load values from .env into the environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}
