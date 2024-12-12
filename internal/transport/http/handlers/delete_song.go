package handlers

import (
	"net/http"
)

func (h *Handlers) DeleteSong(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("DeleteSong: start", "url", r.URL.String())

	id, err := parseIDParam(r)
	if err != nil {
		h.logger.Warn("DeleteSong: missing or invalid id")
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.libraryService.DeleteSong(r.Context(), id); err != nil {
		h.logger.Warn("DeleteSong: failed to delete", "id", id, "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.logger.Info("DeleteSong: success", "id", id)
	w.WriteHeader(http.StatusNoContent)
}
