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

	err := h.service.AddFavoriteMovieToUser(r.Context(), userID, req.MovieID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

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

func (h *FavoriteHandler) DeleteFavoriteMovieFromUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Пользователь не найден", http.StatusUnauthorized)
		return
	}

	movieIDstr := r.PathValue("movie_id")
	movieID, err := strconv.Atoi(movieIDstr)
	if err != nil {
		http.Error(w, "Неверный ID фильма", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteFavoriteMovieFromUser(r.Context(), userID, movieID)
	if err != nil {
		if errors.Is(err, repositories.ErrMovieNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
