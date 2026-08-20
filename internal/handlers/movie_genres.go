package handlers

import (
	"encoding/json"
	"net/http"
	"ozinshe/internal/services"
	"strconv"
)

type MovieGenreHandler struct {
	service *services.MovieGenreService
}

func NewMovieGenreHandler(service *services.MovieGenreService) *MovieGenreHandler {
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

func (h *MovieGenreHandler) GetGenresOfMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	genres, err := h.service.GetGenresOfMovie(ctx, id)
	if err := json.NewEncoder(w).Encode(genres); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func (h *MovieGenreHandler) DeleteGenreFromMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var req genreIDRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteGenreFromMovie(ctx, id, req.GenreID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
