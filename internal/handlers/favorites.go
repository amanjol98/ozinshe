package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"ozinshe/internal/middleware"
	"ozinshe/internal/repositories"
	"ozinshe/internal/services"
	"strconv"
)

type FavoriteHandler struct {
	service *services.FavoriteService
}

func NewFavoriteHandler(service *services.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

type addMovieReq struct {
	MovieID int `json:"movie_id"`
}

// AddFavoriteMovieToUser godoc
// @Summary Добавить фильм в избранные
// @Description Добавляет выбранный фильм в список избранных текущего пользователя.
// @Tags Favorites
// @Accept json
// @Security BearerAuth
// @Param request body addMovieReq true "Данные для добавления фильма"
// @Success 201
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /favorites [post]
func (h *FavoriteHandler) AddFavoriteMovieToUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	var req addMovieReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ввели неверный ID фильма", http.StatusBadRequest)
		return
	}

	if req.MovieID <= 0 {
		http.Error(w, "ID фильма должен быть выше 0", http.StatusBadRequest)
		return
	}

	err := h.service.AddFavoriteMovieToUser(r.Context(), userID, req.MovieID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// GetFavoriteMovies godoc
// @Summary Получить список избранных фильмов
// @Description Возвращает список всех избранных фильмов текущего пользователя
// @Tags Favorites
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.FavoriteMovieResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /favorites [get]
func (h *FavoriteHandler) GetFavoriteMovies(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	favoriteMovieResponse, err := h.service.GetFavoriteMovies(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repositories.ErrFavoriteNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(favoriteMovieResponse); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteFavoriteMovieFromUser godoc
// @Summary Удалить фильм из избранного
// @Description Удаляет выбранный фильм из списка избранных текущего пользователя.
// @Tags Favorites
// @Security BearerAuth
// @Param movie_id path int true "ID фильма"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /favorites/{movie_id} [delete]
func (h *FavoriteHandler) DeleteFavoriteMovieFromUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	movieIDStr := r.PathValue("movie_id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		http.Error(w, "Неверный ID фильма", http.StatusBadRequest)
		return
	}

	if movieID <= 0 {
		http.Error(w, "ID должен быть выше 0", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteFavoriteMovieFromUser(r.Context(), userID, movieID)
	if err != nil {
		if errors.Is(err, repositories.ErrFavoriteNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
