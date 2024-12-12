package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/1abobik1/online_song_lib/internal/config"
	"github.com/1abobik1/online_song_lib/internal/repository"
	"github.com/1abobik1/online_song_lib/internal/service"
	"github.com/1abobik1/online_song_lib/internal/storage/postgresql"
	httpTransport "github.com/1abobik1/online_song_lib/internal/transport/http"
	"github.com/1abobik1/online_song_lib/internal/transport/http/handlers"

	"github.com/babenow/slogwrapper/slogpretty"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// загрузка конфига
	cfg := config.MustLoad()
	// установка логгера
	logger := setupLogger(cfg.Env)

	logger.Info("ready config", slog.Any("config", cfg))

	// создает новое подключение к БД
	storage, err := postgresql.New(cfg.StoragePath, logger)
	if err != nil {
		logger.Error("failed to init storage", "error", err)
		os.Exit(1)
	}

	var repo repository.SongRepository = storage
	libService := service.NewLibraryService(repo)

	// Создание хендлеров из пакета handlers
	h := handlers.NewHandlers(libService, logger, cfg)

	// Создаение роутера
	r := httpTransport.NewRouter(h)

	addr := cfg.HTTPServer
	logger.Info("server listening on " + addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal, envDev:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug},
		}
		handler := opts.NewPrettyHandler(os.Stdout)
		log = slog.New(handler)
	case envProd:
		opts := slogpretty.PrettyHandlerOptions{
			SlogOpts: &slog.HandlerOptions{Level: slog.LevelInfo},
		}
		handler := opts.NewPrettyHandler(os.Stdout)
		log = slog.New(handler)
	}

	return log
}
