package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"ozinshe/internal/repositories"
	"ozinshe/internal/services"
	"strconv"
)

type GenreHandler struct {
	service *services.GenreService
}

func NewGenreHandler(service *services.GenreService) *GenreHandler {
	return &GenreHandler{service: service}
}

type genreRequest struct {
	Name string `json:"name"`
}

// GetAllGenres godoc
// @Summary Получить список всех жанров
// @Description Возвращает списко всех доступных жанров.
// @Tags Genres
// @Produce json
// @Success 200 {array} models.Genre
// @Failure 500 {object} map[string]string
// @Router /genres [get]
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

// GetGenreByID godoc
// @Summary Получить жанр по ID
// @Description Возвращает выбранный по ID жанр
// @Tags Genres
// @Produce json
// @Param id path int true "ID жанра"
// @Success 200 {object} models.Genre
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /genres/{id} [get]
func (h *GenreHandler) GetGenreByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "ID должен быть больше 0", http.StatusBadRequest)
		return
	}

	genre, err := h.service.GetGenreByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrGenreNotFound) {
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

// CreateGenre godoc
// @Summary Создать жанр
// @Description Создает новый жанр. Доступно только администраторам.
// @Tags Genres
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body genreRequest true "Данные для создания"
// @Success 201 {object} models.Genre
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /genres [post]
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

// DeleteGenre godoc
// @Summary Удалить жанр
// @Description Удаляет выбранный жанр. Доступно только администраторам.
// @Tags Genres
// @Security BearerAuth
// @Param id path int true "ID жанра"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /genres/{id} [delete]
func (h *GenreHandler) DeleteGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "ID должен быть больше 0", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteGenre(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrGenreNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Неизвестная ошибка во время удаления жанра", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateGenre godoc
// @Summary Обновить жанр
// @Description Обновляет выбранный жанр. Доступно только администраторам.
// @Tags Genres
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID жанра"
// @Param request body genreRequest true "Данные для обновления"
// @Success 200 {object} models.Genre
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /genres/{id} [patch]
func (h *GenreHandler) UpdateGenre(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request genreRequest

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "ID должен быть больше 0", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Ошибка в JSON", http.StatusBadRequest)
		return
	}

	genre, err := h.service.UpdateGenre(ctx, request.Name, id)
	if err != nil {
		if errors.Is(err, repositories.ErrGenreNotFound) {
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
