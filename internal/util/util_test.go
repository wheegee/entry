package util

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		expectedLevel zerolog.Level
	}{
		{name: "panic level", envValue: "panic", expectedLevel: zerolog.PanicLevel},
		{name: "fatal level", envValue: "fatal", expectedLevel: zerolog.FatalLevel},
		{name: "error level", envValue: "error", expectedLevel: zerolog.ErrorLevel},
		{name: "warn level", envValue: "warn", expectedLevel: zerolog.WarnLevel},
		{name: "info level", envValue: "info", expectedLevel: zerolog.InfoLevel},
		{name: "debug level", envValue: "debug", expectedLevel: zerolog.DebugLevel},
		{name: "trace level", envValue: "trace", expectedLevel: zerolog.TraceLevel},
		{name: "default level for invalid input", envValue: "invalid", expectedLevel: zerolog.WarnLevel},
		{name: "default level when no env var set", envValue: "", expectedLevel: zerolog.WarnLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the environment variable if provided
			if tt.envValue != "" {
				os.Setenv("LOG_LEVEL", tt.envValue)
			} else {
				os.Unsetenv("LOG_LEVEL")
			}

			// Call the function
			SetLogLevel()

			// Check the global log level
			assert.Equal(t, tt.expectedLevel, zerolog.GlobalLevel())

			// Clean up the environment variable for next test
			os.Unsetenv("LOG_LEVEL")
		})
	}
}
