package log

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
}

// Info returns a new event with info level.
func Info() *zerolog.Event {
	return log.Logger.Info()
}

// Error returns a new event with error level.
func Error() *zerolog.Event {
	return log.Logger.Error()
}

// Warn returns a new event with warn level.
func Warn() *zerolog.Event {
	return log.Logger.Warn()
}

// Fatal returns a new event with fatal level.
func Fatal() *zerolog.Event {
	return log.Logger.Fatal()
}
