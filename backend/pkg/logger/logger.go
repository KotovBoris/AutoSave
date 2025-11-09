package logger

import (
    "os"
    "time"

    "github.com/rs/zerolog"
)

type Logger struct {
    *zerolog.Logger
}

func New(level zerolog.Level, format string) *Logger {
    var logger zerolog.Logger

    if format == "json" {
        logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
    } else {
        output := zerolog.ConsoleWriter{
            Out:        os.Stdout,
            TimeFormat: time.RFC3339,
            NoColor:    false,
        }
        logger = zerolog.New(output).With().Timestamp().Logger()
    }

    logger = logger.Level(level)

    return &Logger{&logger}
}

func NewDefault() *Logger {
    return New(zerolog.DebugLevel, "console")
}

func (l *Logger) WithField(key string, value interface{}) *Logger {
    logger := l.Logger.With().Interface(key, value).Logger()
    return &Logger{&logger}
}

func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
    logger := l.Logger.With().Fields(fields).Logger()
    return &Logger{&logger}
}

func (l *Logger) WithError(err error) *Logger {
    logger := l.Logger.With().Err(err).Logger()
    return &Logger{&logger}
}

