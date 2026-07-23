package genres

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

	genres, err := h.service.GetAll(ctx)
	if err != nil {
		http.Error(w, "Не удалось получить список жанров", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(genres); err != nil {
		http.Error(w, "Ошибка при загрузке жанров в JSON", http.StatusInternalServerError)
		return
	}
}

func (h *GenreHandler) GetGenreByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	genre, err := h.service.GetGenreByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrGenreNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Не удалось получить жанр", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(genre); err != nil {
		http.Error(w, "Ошибка при загрузке жанра в JSON", http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(genre); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusInternalServerError)
		return
	}
}

func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteGenre(ctx, id)
	if err != nil {
		if errors.Is(err, ErrGenreNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Неизвестная ошибка во время удаления жанра", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GenreHandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request genreRequest

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusBadRequest)
		return
	}

	genre, err := h.service.UpdateGenre(ctx, request.Name, id)
	if err != nil {
		if errors.Is(err, ErrGenreNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(genre); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusInternalServerError)
		return
	}
}
