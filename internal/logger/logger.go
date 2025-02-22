package logger

import (
	"log/slog"
	"os"
)

func Setup(env string) (logger *slog.Logger) {
	if env == "local" {
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	} else {
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	slog.SetDefault(logger)

	return
}
