package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string
	GRPCPort     string
	MetadataAddr string
	RedisAddr    string
	UploadDir    string
}

func Load() Config {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	return Config{
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		GRPCPort:     getEnv("GRPC_PORT", "50051"),
		MetadataAddr: getEnv("METADATA_SERVICE_ADDR", "localhost:50051"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		UploadDir:    getEnv("UPLOAD_DIR", "./data/uploads"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
