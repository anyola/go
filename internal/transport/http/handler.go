package httptransport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"movieapi/internal/movie"
	"movieapi/internal/service"
)

type Handler struct{ service *service.Service }
func NewHandler(s *service.Service) *Handler { return &Handler{service: s} }

func (h *Handler) LoadMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method must be GET")
		return
	}
	from, to, err := parseRange(r)
	if err != nil { writeError(w, http.StatusBadRequest, err.Error()); return }
	movies, errs := h.service.LoadMovies(r.Context(), from, to)
	if len(errs) == 0 {
		writeJSON(w, http.StatusOK, movies)
		return
	}
	writeJSON(w, http.StatusOK, struct {
		Movies []*movie.Movie    `json:"movies"`
		Errors []movie.LoadError `json:"errors"`
	}{movies, errs})
}
func (h *Handler) LastMovie(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method must be GET")
		return
	}
	m, err := h.service.LastMovie(r.Context())
	if err != nil {
		status := http.StatusBadGateway
		if errors.Is(err, movie.ErrNotFound) { status = http.StatusNotFound }
		if errors.Is(err, context.DeadlineExceeded) { status = http.StatusGatewayTimeout }
		writeError(w, status, err.Error()); return
	}
	writeJSON(w, http.StatusOK, m)
}
func parseRange(r *http.Request) (int, int, error) {
	from, err := strconv.Atoi(r.URL.Query().Get("from"))
	if err != nil || from <= 0 { return 0, 0, errors.New("from must be a positive integer") }
	to, err := strconv.Atoi(r.URL.Query().Get("to"))
	if err != nil || to <= 0 { return 0, 0, errors.New("to must be a positive integer") }
	if from > to { return 0, 0, errors.New("from must not be greater than to") }
	if to-from > 10000 { return 0, 0, errors.New("range must not exceed 10000 movies") }
	return from, to, nil
}
func writeJSON(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
