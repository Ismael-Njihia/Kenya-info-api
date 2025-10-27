package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	URI     string
	Name    string
	Timeout int
}

type LoggerConfig struct {
	Level string
}

func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load(".env")

	// Parse timeout from string to int
	timeout, err := strconv.Atoi(getEnv("MONGODB_TIMEOUT", "10"))
	if err != nil {
		log.Printf("Invalid MONGODB_TIMEOUT, using default 10s")
		timeout = 10
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			URI:     getEnv("MONGODB_URI", ""),
			Name:    getEnv("MONGODB_DATABASE", "kenya_info"),
			Timeout: timeout,
		},
		Logger: LoggerConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
