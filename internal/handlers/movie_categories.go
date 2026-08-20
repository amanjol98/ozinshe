package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"ozinshe/internal/repositories"
	"ozinshe/internal/services"
	"strconv"
)

type MovieCategoryHandler struct {
	service *services.MovieCategoryService
}

func NewMovieCategoryHandler(service *services.MovieCategoryService) *MovieCategoryHandler {
	return &MovieCategoryHandler{service: service}
}

type categoryRequest struct {
	CategoryID int `json:"category_id"`
}

func (h *MovieCategoryHandler) AddCategoryToMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req categoryRequest

	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.AddCategoryToMovie(ctx, id, req.CategoryID)
	if err != nil {
		http.Error(w, "Ошибка при добавлении категории", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *MovieCategoryHandler) GetMovieCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	categories, err := h.service.GetMovieCategories(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrMovieCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(categories); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MovieCategoryHandler) DeleteCategoryFromMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req categoryRequest

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCategoryFromMovie(ctx, id, req.CategoryID)
	if err != nil {
		if errors.Is(err, repositories.ErrMovieCategoryNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
