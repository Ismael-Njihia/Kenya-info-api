package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad_WithDefaults(t *testing.T) {
	// Clear environment variables
	os.Clearenv()
	os.Setenv("MONGODB_URI", "mongodb://localhost:27017")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "release", cfg.Server.Mode)
	assert.Equal(t, "kenya_info", cfg.Database.Database)
}

func TestLoad_WithCustomValues(t *testing.T) {
	os.Clearenv()
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("SERVER_MODE", "debug")
	os.Setenv("MONGODB_URI", "mongodb://testhost:27017")
	os.Setenv("MONGODB_DATABASE", "test_db")
	os.Setenv("LOG_LEVEL", "debug")

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "9090", cfg.Server.Port)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, "debug", cfg.Server.Mode)
	assert.Equal(t, "test_db", cfg.Database.Database)
	assert.Equal(t, "debug", cfg.Logger.Level)
}

func TestLoad_MissingURI(t *testing.T) {
	os.Clearenv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	// When MONGODB_URI is not set, it defaults to mongodb://localhost:27017
	assert.Equal(t, "mongodb://localhost:27017", cfg.Database.URI)
}
