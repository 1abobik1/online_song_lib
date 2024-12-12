package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// parseIntQueryParam пытается взять целочисленное значение из query параметра.
func parseIntQueryParam(r *http.Request, key string) (int, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return 0, nil // Параметра нет — вернем 0
	}
	return strconv.Atoi(valStr)
}

// parseIDParam извлекает и парсит ID из URL-параметра.
func parseIDParam(r *http.Request) (int, error) {
	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		return 0, http.ErrNoLocation
	}
	return strconv.Atoi(idStr)
}

// writeJSON упрощает запись JSON-ответа.
func writeJSON(w http.ResponseWriter, v interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writePlainText упрощает запись текстового ответа.
func writePlainText(w http.ResponseWriter, text string, status int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(text))
}
