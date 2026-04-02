package logging

import (
	"log/slog"
	"os"
)

type Log struct {
	Slogger *slog.Logger
	opts    *slog.HandlerOptions
}

func New() *Log {
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	}

	slogger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	return &Log{
		opts:    opts,
		Slogger: slogger,
	}
}

func NewLoggerForTest() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
