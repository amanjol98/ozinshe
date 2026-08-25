package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"ozinshe/internal/repositories"
	"ozinshe/internal/services"
	"strconv"
)

type MovieScreenshotsHandler struct {
	service *services.MovieScreenshotsService
}

func NewMovieScreenshotsHandler(service *services.MovieScreenshotsService) *MovieScreenshotsHandler {
	return &MovieScreenshotsHandler{service: service}
}

type screenshotRequest struct {
	ImageURL string `json:"image_url"`
}

func (h *MovieScreenshotsHandler) AddScreenshotsToMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req screenshotRequest

	idString := r.PathValue("id")
	movieID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.AddScreenshotsToMovie(ctx, movieID, req.ImageURL)
	if err != nil {
		http.Error(w, "Ошибка при добавлении скриншота к фильму", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (h *MovieScreenshotsHandler) GetScreenshots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	movieID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	movieScreenshots, err := h.service.GetScreenshots(ctx, movieID)
	if err != nil {
		http.Error(w, "Ошибка при получении скриншота", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(movieScreenshots); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *MovieScreenshotsHandler) DeleteScreenshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")
	screenshotID, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteScreenshot(ctx, screenshotID)
	if err != nil {
		if errors.Is(err, repositories.ErrScreenshotNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
