package logger

import (
	"log/slog"
	"os"
)

func Init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	lgr := slog.New(handler)

	slog.SetDefault(lgr)
}
