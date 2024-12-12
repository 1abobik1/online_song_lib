package handlers

import (
	"net/http"

	"github.com/1abobik1/online_song_lib/internal/repository"
)

// GetSongs - обработчик для GET /library
func (h *Handlers) GetSongs(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GetSongs: start", "url", r.URL.String())

	limit, err := parseIntQueryParam(r, "limit")
	if err != nil {
		h.logger.Warn("GetSongs: invalid limit")
		http.Error(w, "invalid limit", http.StatusBadRequest)
		return
	}

	offset, err := parseIntQueryParam(r, "offset")
	if err != nil {
		h.logger.Warn("GetSongs: invalid offset")
		http.Error(w, "invalid offset", http.StatusBadRequest)
		return
	}

	filter := repository.SongFilter{
		GroupName:   r.URL.Query().Get("group"),
		SongName:    r.URL.Query().Get("song"),
		ReleaseDate: r.URL.Query().Get("releaseDate"),
		Limit:       limit,
		Offset:      offset,
	}

	songs, err := h.libraryService.GetSongs(r.Context(), filter)
	if err != nil {
		h.logger.Warn("GetSongs: failed to get songs", "error", err)
		http.Error(w, "failed to get songs", http.StatusInternalServerError)
		return
	}

	h.logger.Info("GetSongs: success", "count", len(songs))
	writeJSON(w, songs, http.StatusOK)
}
