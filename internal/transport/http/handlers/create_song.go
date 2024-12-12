package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/1abobik1/online_song_lib/internal/transport/http/dto"
)

// POST /library
func (h *Handlers) CreateSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("CreateSong: start", "url", r.URL.String())

	var req dto.CreateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("CreateSong: invalid request body")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Group == "" || req.Song == "" {
		h.logger.Warn("CreateSong: group or song is empty")
		http.Error(w, "group and song fields are required", http.StatusBadRequest)
		return
	}

	// обращение к внешнему API
	externalAPIURL := h.cfg.ExternalApiURL

	newSong, err := h.libraryService.AddSong(r.Context(), req.Group, req.Song, externalAPIURL)
	if err != nil {
		h.logger.Warn("CreateSong: failed to add song", "error", err)
		http.Error(w, "failed to add song", http.StatusInternalServerError)
		return
	}

	h.logger.Info("CreateSong: success", "id", newSong.ID)
	writeJSON(w, newSong, http.StatusCreated)
}
