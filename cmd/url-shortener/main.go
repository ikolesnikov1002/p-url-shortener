package main

import (
	"log/slog"
	"os"
	"url-shortener/internal/config"
	sl "url-shortener/internal/lib/logger"
	"url-shortener/internal/storage/sqlite"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Set up logger
	log := setupLogger(cfg.Env)

	// Init storage
	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
	}

	log.Info("Starting server", slog.String("Address", cfg.Host+":"+cfg.Port))
	log.Debug("Debug mode is enabled.")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "production":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
