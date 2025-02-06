package logger

import (
	"context"
	"encoding/json"
	"io"
	"os"

	"log"
	"log/slog"

	"github.com/fatih/color"
)

const (
	DEV = 1 + iota
	PROD
)

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := color.CyanString(r.Message)

	h.l.Println(timeStr, level, msg, color.WhiteString(string(b)))

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}
func SetupLogger(
	appEnv int,
	filepath string,
) *slog.Logger {
	var logger *slog.Logger
	logger = slog.Default()
	if appEnv == DEV {
		opts := PrettyHandlerOptions{
			SlogOpts: slog.HandlerOptions{
				Level: slog.LevelDebug,
			},
		}
		handler := NewPrettyHandler(os.Stdout, opts)
		logger = slog.New(handler)
		return logger
	}
	if appEnv == PROD {
		opts := &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
		var handler *slog.JSONHandler
		f, err := os.Open(filepath)
		if err != nil {
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			handler = slog.NewJSONHandler(f, opts)
		}
		logger = slog.New(handler)
		return logger
	}
	return logger
}
