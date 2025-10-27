package config

import (
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Logger   LoggerConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
	Mode string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URI      string
	Database string
	Timeout  int
}

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Level      string
	OutputPath string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Mode: getEnv("SERVER_MODE", "release"),
		},
		Database: DatabaseConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "kenya_info"),
			Timeout:  getEnvAsInt("MONGODB_TIMEOUT", 10),
		},
		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			OutputPath: getEnv("LOG_OUTPUT_PATH", "stdout"),
		},
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
