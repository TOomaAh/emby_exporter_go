package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type Interface interface {
	Debug(format interface{}, args ...interface{})
	Info(format string, args ...interface{})
	Warn(format string, args ...interface{})
	Error(format interface{}, args ...interface{})
	Fatal(format interface{}, args ...interface{})
}

// Logger -.
type Logger struct {
	logger *zerolog.Logger
}

var _ Interface = (*Logger)(nil)

func New(level string) *Logger {

	var l zerolog.Level

	switch strings.ToLower(level) {
	case "error":
		l = zerolog.ErrorLevel
	case "warn":
		l = zerolog.WarnLevel
	case "info":
		l = zerolog.InfoLevel
	case "debug":
		l = zerolog.DebugLevel
	default:
		l = zerolog.InfoLevel
	}

	skipFrameCount := 2

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}

	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}

	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}

	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}

	logger := zerolog.New(output).Level(l).With().Timestamp().CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + skipFrameCount).Logger()

	return &Logger{
		logger: &logger,
	}

}

// Debug -.
func (l *Logger) Debug(message interface{}, args ...interface{}) {
	l.msg(l.logger.Debug(), message, args...)
}

// Info -.
func (l *Logger) Info(message string, args ...interface{}) {
	l.log(l.logger.Info(), message, args...)
}

// Warn -.
func (l *Logger) Warn(message string, args ...interface{}) {
	l.log(l.logger.Warn(), message, args...)
}

// Error -.
func (l *Logger) Error(message interface{}, args ...interface{}) {
	l.msg(l.logger.Error(), message, args...)
}

// Fatal -.
func (l *Logger) Fatal(message interface{}, args ...interface{}) {
	l.msg(l.logger.Fatal(), message, args...)
}

func (l *Logger) log(event *zerolog.Event, message string, args ...interface{}) {
	if len(args) == 0 {
		event.Msg(message)
	} else {
		event.Msgf(message, args...)
	}
}

func (l *Logger) msg(event *zerolog.Event, message interface{}, args ...interface{}) {
	var text string
	switch m := message.(type) {
	case error:
		text = m.Error()
	case string:
		text = m
	default:
		text = fmt.Sprintf("%v", m)
	}

	if len(args) == 0 {
		event.Msg(text)
	} else {
		event.Msgf(text, args...)
	}
}
