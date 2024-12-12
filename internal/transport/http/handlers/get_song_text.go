package handlers

import (
	"net/http"
)

func (h *Handlers) GetSongText(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("GetSongText: start", "url", r.URL.String())

	id, err := parseIDParam(r)
	if err != nil {
		h.logger.Warn("GetSongText: missing or invalid id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	verse, err := parseIntQueryParam(r, "verse")
	if err != nil {
		h.logger.Warn("GetSongText: invalid verse")
		http.Error(w, "invalid verse", http.StatusBadRequest)
		return
	}

	text, err := h.libraryService.GetSongTextByVerse(r.Context(), id, verse)
	if err != nil {
		h.logger.Warn("GetSongText: not found", "id", id, "verse", verse)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info("GetSongText: success", "id", id)
	writePlainText(w, text, http.StatusOK)
}
