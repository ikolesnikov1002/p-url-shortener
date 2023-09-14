package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	//rec, err := storage.CreateUrl("http://google.com/3")
	//
	//if err != nil {
	//	log.Error("Failed to create url", sl.Err(err))
	//}
	//log.Info("saved id", slog.Int64("ID", rec))

	//res, err := storage.GetUrl("B8t3gG")
	//if err != nil {
	//	log.Error("Failed to get url", sl.Err(err))
	//}
	//log.Info("saved url", slog.String("URL", res))

	//err = storage.DeleteUrl("lLcLzH")
	//if err != nil {
	//	log.Error("Failed to delete url", sl.Err(err))
	//}
	//log.Info("url deleted")

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
