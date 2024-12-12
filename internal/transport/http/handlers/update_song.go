package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/1abobik1/online_song_lib/internal/repository"
	"github.com/1abobik1/online_song_lib/internal/transport/http/dto"
)

func (h *Handlers) UpdateSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("UpdateSong: start", "url", r.URL.String())

	id, err := parseIDParam(r)
	if err != nil {
		h.logger.Warn("UpdateSong: missing or invalid id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("UpdateSong: invalid request body")
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	update := repository.SongUpdate{
		GroupName:   req.GroupName,
		SongName:    req.SongName,
		ReleaseDate: req.ReleaseDate,
		Text:        req.Text,
		Link:        req.Link,
	}

	if err := h.libraryService.UpdateSong(r.Context(), id, update); err != nil {
		h.logger.Warn("UpdateSong: failed to update", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info("UpdateSong: success", "id", id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("song updated successfully"))
}
