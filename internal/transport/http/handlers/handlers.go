package handlers

import (
	"log/slog"

	"github.com/1abobik1/online_song_lib/internal/config"
	"github.com/1abobik1/online_song_lib/internal/service"
)

type Handlers struct {
	libraryService *service.LibraryService
	logger         *slog.Logger
	cfg            *config.Config
}

func NewHandlers(libraryService *service.LibraryService, logger *slog.Logger, cfg *config.Config) *Handlers {
    return &Handlers{
        libraryService: libraryService,
        logger: logger,
        cfg: cfg,
    }
}

