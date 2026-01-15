package logger_test

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"go-boilerplate-rest-api-chi/internal/config"
	"go-boilerplate-rest-api-chi/internal/logger"
)

func TestNewLogger(t *testing.T) {
	t.Run("valid debug level", func(t *testing.T) {
		cfg := &config.Config{Log: config.LogConfig{Level: "debug", Format: "json"}}
		log, err := logger.NewLogger(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, log)
		assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
	})

	t.Run("invalid level defaults to info", func(t *testing.T) {
		cfg := &config.Config{Log: config.LogConfig{Level: "unknown", Format: "json"}}
		log, err := logger.NewLogger(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, log)
		assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
	})

	t.Run("text format creates console writer", func(t *testing.T) {
		cfg := &config.Config{Log: config.LogConfig{Level: "info", Format: "text"}}
		log, err := logger.NewLogger(cfg)

		assert.NoError(t, err)
		assert.NotNil(t, log)
	})
}
