package movie_genres

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type MovieGenreHandler struct {
	service *MovieGenreService
}

func NewMovieGenreHandler(service *MovieGenreService) *MovieGenreHandler {
	return &MovieGenreHandler{service: service}
}

type genreIDRequest struct {
	GenreID int `json:"genre_id"`
}

func (h *MovieGenreHandler) AddGenreToMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	movieID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var req genreIDRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.AddGenreToMovie(ctx, movieID, req.GenreID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
