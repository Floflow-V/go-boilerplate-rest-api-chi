package logger

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"

	"go-boilerplate-rest-api-chi/internal/config"
)

func NewLogger(cfg *config.Config) (zerolog.Logger, error) {
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(level)

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	logger := zerolog.New(os.Stderr).With().Timestamp().Caller().Logger()

	loc, err := time.LoadLocation("Europe/Paris")
	if err != nil {
		loc = time.Local
	}

	if cfg.Log.Format == "text" {

		logger = logger.Output(zerolog.ConsoleWriter{
			Out:          os.Stderr,
			TimeLocation: loc,
			TimeFormat:   "15:04:05 02/01/2006",
		})
	}

	return logger, nil
}
