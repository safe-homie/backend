package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

type logger struct {
	log zerolog.Logger
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
}

func New() Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("*%s*", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}
	return &logger{
		log: zerolog.New(output).With().Timestamp().Logger(),
	}
}

func (l *logger) Info(msg string) {
	l.log.Info().Msg(msg)
}

func (l *logger) Error(msg string) {
	l.log.Error().Msg(msg)
}

func (l *logger) Debug(msg string) {
	l.log.Debug().Msg(msg)
}
