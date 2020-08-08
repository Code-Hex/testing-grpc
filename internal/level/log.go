package level

import (
	"strings"

	"github.com/rs/zerolog"
)

func Log(s string) zerolog.Level {
	switch strings.ToUpper(s) {
	case "INFO":
		return zerolog.InfoLevel
	case "WARN":
		return zerolog.WarnLevel
	case "ERROR":
		return zerolog.ErrorLevel
	}
	return zerolog.DebugLevel
}
