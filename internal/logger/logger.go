package logger

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	opts := &Options{Level: slog.LevelInfo, TimeFormat: "01-02|15:04:05.000"}
	handler := NewHandler(os.Stdout, opts)
	return slog.New(handler)
}
