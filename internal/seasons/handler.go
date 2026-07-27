package seasons

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type SeasonHandler struct {
	service *SeasonService
}

func NewSeasonHandler(service *SeasonService) *SeasonHandler {
	return &SeasonHandler{service: service}
}

type seasonRequest struct {
	SeasonNumber int `json:"season_number"`
}

func (h *SeasonHandler) AddSeasonToMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req seasonRequest

	idString := r.PathValue("movieID")

	movieID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.AddSeasonToMovie(ctx, movieID, req.SeasonNumber)
	if err != nil {
		if errors.Is(err, ErrSeasonNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка при добавлении сезона в фильм", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *SeasonHandler) DeleteSeasonFromMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	sID := r.PathValue("seasonID")

	seasonID, err := strconv.Atoi(sID)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteSeasonFromMovie(ctx, seasonID)

	if err != nil {
		if errors.Is(err, ErrSeasonNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
