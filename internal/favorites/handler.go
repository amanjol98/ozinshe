package favorites

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type FavoriteHandler struct {
	service *FavoriteService
}

func NewFavoriteHandler(service *FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

type addMovieReq struct {
	MovieID int `json:"movie_id"`
}

func (h *FavoriteHandler) AddFavoriteMovieToUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req addMovieReq

	idString := r.PathValue("user_id")
	userID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Ввели неверный ID фильма", http.StatusBadRequest)
		return
	}

	err = h.service.AddFavoriteMovieToUser(ctx, userID, req.MovieID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
