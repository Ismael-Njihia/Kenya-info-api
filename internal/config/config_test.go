package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadDefaultConfig(t *testing.T) {
	// Set required env var
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")
	defer os.Unsetenv("MONGODB_URI")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "development", cfg.Server.Env)
	assert.Equal(t, "kenya_info", cfg.Database.Database)
	assert.Equal(t, 10, cfg.Database.Timeout)
	assert.Equal(t, "info", cfg.Logger.Level)
}

func TestLoadCustomConfig(t *testing.T) {
	os.Setenv("PORT", "3000")
	os.Setenv("ENV", "production")
	os.Setenv("MONGODB_URI", "mongodb://custom:27017")
	os.Setenv("MONGODB_DATABASE", "custom_db")
	os.Setenv("MONGODB_TIMEOUT", "30")
	os.Setenv("LOG_LEVEL", "debug")

	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("ENV")
		os.Unsetenv("MONGODB_URI")
		os.Unsetenv("MONGODB_DATABASE")
		os.Unsetenv("MONGODB_TIMEOUT")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "3000", cfg.Server.Port)
	assert.Equal(t, "production", cfg.Server.Env)
	assert.Equal(t, "mongodb://custom:27017", cfg.Database.MongoURI)
	assert.Equal(t, "custom_db", cfg.Database.Database)
	assert.Equal(t, 30, cfg.Database.Timeout)
	assert.Equal(t, "debug", cfg.Logger.Level)
}

func TestLoadConfigMissingMongoURI(t *testing.T) {
	// Store and clear the MONGODB_URI env var
	originalURI := os.Getenv("MONGODB_URI")
	os.Unsetenv("MONGODB_URI")
	
	// Restore it after test
	defer func() {
		if originalURI != "" {
			os.Setenv("MONGODB_URI", originalURI)
		}
	}()

	cfg, err := Load()
	assert.Error(t, err)
	assert.Nil(t, cfg)
	if err != nil {
		assert.Contains(t, err.Error(), "MONGODB_URI is required")
	}
}

func TestGetEnv(t *testing.T) {
	os.Setenv("TEST_VAR", "test_value")
	defer os.Unsetenv("TEST_VAR")

	value := getEnv("TEST_VAR", "default")
	assert.Equal(t, "test_value", value)

	value = getEnv("NON_EXISTENT", "default")
	assert.Equal(t, "default", value)
}

func TestGetEnvAsInt(t *testing.T) {
	os.Setenv("TEST_INT", "42")
	defer os.Unsetenv("TEST_INT")

	value := getEnvAsInt("TEST_INT", 10)
	assert.Equal(t, 42, value)

	value = getEnvAsInt("NON_EXISTENT", 10)
	assert.Equal(t, 10, value)

	os.Setenv("INVALID_INT", "not_a_number")
	defer os.Unsetenv("INVALID_INT")

	value = getEnvAsInt("INVALID_INT", 10)
	assert.Equal(t, 10, value)
}
