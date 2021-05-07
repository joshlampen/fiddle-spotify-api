package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() (zerolog.Logger) {
    consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
    logger := zerolog.New(consoleWriter).With().Timestamp().Logger()

    return logger
}
