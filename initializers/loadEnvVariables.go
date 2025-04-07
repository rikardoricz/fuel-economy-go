package initializers

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadEnvVariables() {
	if os.Getenv("GO_ENV") == "production" {
		log.Println("Running in production mode, using environment variables")
		return
	}

	// for development try to load from .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, will use environment variables")
	}
}
