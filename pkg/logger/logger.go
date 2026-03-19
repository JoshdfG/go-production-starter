package logger

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func New(level string, appName string, version string, env string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.CallerMarshalFunc = shortCaller

	var output io.Writer
	if env == "development" {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "15:04:05",
		}
	} else {
		output = os.Stdout
	}

	logger := zerolog.New(output).
		With().
		Timestamp().
		Caller().
		Str("app", appName).
		Str("version", version).
		Logger()

	log.Logger = logger
	return logger
}

// shortCaller trims the full path to just "file.go:line"
// using fmt.Sprintf directly — never calls zerolog.CallerMarshalFunc
// to avoid infinite recursion
func shortCaller(_ uintptr, file string, line int) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	return fmt.Sprintf("%s:%d", short, line)
}
