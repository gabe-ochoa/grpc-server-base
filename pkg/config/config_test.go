package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustLoadConfig_defaults(t *testing.T) {
	cfg := MustLoadConfig()

	assert.Equal(t, cfg.ServerMode, Development)
	assert.Equal(t, cfg.LogLevel, "debug")
	assert.Equal(t, cfg.LogFormat, "text")
}

func TestMustLoadConfig_env_override(t *testing.T) {
	os.Setenv("LOG_LEVEL", "info")
	cfg := MustLoadConfig()

	assert.Equal(t, cfg.LogLevel, "info")
}
