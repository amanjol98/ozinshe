package genres

import (
	"encoding/json"
	"net/http"
)

type GenreHandler struct {
	service *GenreService
}

func NewGenreHandler(service *GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

type genreRequest struct {
	Name string `json:"name"`
}

func (h *GenreHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	w.Header().Set("Content-Type", "application/json")

	genre, err := h.service.GetAll(ctx)
	if err != nil {
		http.Error(w, "Не удалось получить список жанров", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(genre); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusInternalServerError)
		return
	}
}

func (h *GenreHandler) CreateGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request genreRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	genre, err := h.service.CreateGenre(ctx, request.Name)
	if err != nil {
		http.Error(w, "Ошибка в создании жанра", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(genre); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusInternalServerError)
		return
	}
}
