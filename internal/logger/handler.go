package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

const bufSize = 1024

type Options struct {
	Level      slog.Leveler
	TimeFormat string
}

type Handler struct {
	opts Options
	mu   *sync.Mutex
	out  io.Writer
}

func NewHandler(out io.Writer, opts *Options) *Handler {
	h := &Handler{out: out, mu: new(sync.Mutex)}
	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	if h.opts.TimeFormat == "" {
		h.opts.TimeFormat = time.RFC3339
	}
	return h
}

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *Handler) Handle(_ context.Context, record slog.Record) error {
	buf := make([]byte, 0, bufSize)

	// time
	if !record.Time.IsZero() {
		buf = fmt.Appendf(buf, "[%s] ", record.Time.Format(h.opts.TimeFormat))
	}

	// level
	var level string
	switch record.Level {
	case slog.LevelDebug:
		level = "DBG"
	case slog.LevelInfo:
		level = "INF"
	case slog.LevelWarn:
		level = "WRN"
	case slog.LevelError:
		level = "ERR"
	}
	buf = fmt.Appendf(buf, "%s ", level)

	// message
	buf = fmt.Appendf(buf, "%q ", record.Message)

	// source
	if record.PC != 0 {
		frame, _ := runtime.CallersFrames([]uintptr{record.PC}).Next()
		if frame.PC != 0 {
			buf = fmt.Appendf(buf, "%s:%d", filepath.Base(frame.File), frame.Line)
		}
	}

	// attributes
	record.Attrs(func(a slog.Attr) bool {
		a.Value = a.Value.Resolve()
		if !a.Equal(slog.Attr{}) {
			buf = fmt.Appendf(buf, " %s=%q", a.Key, a.Value.String())
		}
		return true
	})

	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	return err
}

func (h *Handler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *Handler) WithGroup(_ string) slog.Handler {
	return h
}
