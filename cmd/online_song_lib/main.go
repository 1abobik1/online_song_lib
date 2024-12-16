package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/1abobik1/online_song_lib/internal/config"
	"github.com/1abobik1/online_song_lib/internal/repository"
	"github.com/1abobik1/online_song_lib/internal/service"
	"github.com/1abobik1/online_song_lib/internal/storage/postgresql"
	"github.com/1abobik1/online_song_lib/internal/transport/http/handlers"
	"github.com/1abobik1/online_song_lib/internal/transport/http/router"
	"github.com/1abobik1/online_song_lib/pkg"

)


func main() {
	// загрузка конфига
	cfg := config.MustLoad()
	// установка логгера
	logger := pkg.SetupLogger(cfg.Env)

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
	r := router.NewRouter(h)

	addr := cfg.HTTPServer
	logger.Info("server listening on " + addr)
	if err := http.ListenAndServe(addr, pkg.EnableCORS(r)); err != nil {
		logger.Error("server error", "error", err)
		os.Exit(1)
	}
}
