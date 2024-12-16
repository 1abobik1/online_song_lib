package router

import (
	"net/http"

	"github.com/1abobik1/online_song_lib/internal/transport/http/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter(h *handlers.Handlers) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/library", h.GetSongs)
	r.Get("/library/{id}/text", h.GetSongText)
	r.Delete("/library/{id}", h.DeleteSong)
	r.Put("/library/{id}", h.UpdateSong)
	r.Post("/library", h.CreateSong)

	return r
}
