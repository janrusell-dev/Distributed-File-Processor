package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	GRPCPort    string
}

func Load() Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is required")
	}

	port := os.Getenv("GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	return Config{
		DatabaseURL: dbURL,
		GRPCPort:    port,
	}
}
