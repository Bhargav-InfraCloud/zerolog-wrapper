package zerologwrapper

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
)

func init() {
	// Configure the pre-defined field names and options

	// Field name for the text in Msg()
	zerolog.MessageFieldName = "log-message"

	// Field name for timestamp
	zerolog.TimestampFieldName = "timestamp"

	// Field name for log level
	zerolog.LevelFieldName = "log-level"

	// Print just filename and line number for caller field
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return fmt.Sprintf("%s:%d", file[strings.LastIndex(file, "/")+1:], line)
	}
}

// level is a wrapper on top of zerolog's level
type level zerolog.Level

// raw returns the underlying zerolog's level
func (l level) raw() zerolog.Level {
	return zerolog.Level(l)
}

const (
	// Localize log levels from zerolog
	LevelDebug = level(zerolog.DebugLevel)
	LevelInfo  = level(zerolog.InfoLevel)
	LevelWarn  = level(zerolog.WarnLevel)
	LevelError = level(zerolog.ErrorLevel)
	LevelFatal = level(zerolog.FatalLevel)

	// Default log level
	defaultLevel = LevelInfo
)

// Logger has selective logger method signatures
type Logger interface {
	Debug() *zerolog.Event
	Info() *zerolog.Event
	Warn() *zerolog.Event
	Error() *zerolog.Event
	Fatal() *zerolog.Event
	With() zerolog.Context
}

// key holds the type that is used as key to fetch value from context.
// The value in ctx will be a logger so the same logger instance
// can be retrived anywhere the context is available.
type key struct{}

// NewLogger creates a new logger with specified output stream and level.
// This returns the updated context which has key struct as key and logger as value.
// The same logger instance can be retrived from this context.
func NewLogger(ctx context.Context, out io.Writer, level level) (context.Context, Logger) {
	logger := newLogger(out, level)
	ctx = context.WithValue(ctx, key{}, logger)

	return ctx, logger
}

// FromContext retrives the logger from context.
// In absence of logger in context, this creates a new logger with Stdout
// as output stream and default log level specified in constants.
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(key{}).(*zerolog.Logger); ok {
		return logger
	}

	return newLogger(os.Stdout, defaultLevel)
}

// newLogger internally creates a new logger instance enabling
// timestamp, caller info and specified level
func newLogger(out io.Writer, level level) Logger {
	logger := zerolog.New(out).With().
		Timestamp().Caller().Logger().
		Level(level.raw())

	return &logger
}

// FromRawLogger creates Logger type from zerolog Logger
func FromRawLogger(l zerolog.Logger) Logger {
	return &l
}
