package util

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func SetLogLevel() {
	if level, exists := os.LookupEnv("LOG_LEVEL"); exists {
		level = strings.ToLower(level)
		switch level {
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "warn":
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		case "info":
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "trace":
			zerolog.SetGlobalLevel(zerolog.TraceLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		}
		return
	}

	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}
