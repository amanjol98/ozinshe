package movie

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	service *MovieService
}

type movieRequest struct {
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_year"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	PosterURL   string `json:"poster_url"`
	Director    string `json:"director"`
	Producer    string `json:"producer"`
	VideoID     string `json:"video_id"`
}

func NewMovieHandler(service *MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

func (h *MovieHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	movies, err := h.service.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		http.Error(w, "ошибка в JSON:", http.StatusInternalServerError)
		return
	}
}

func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	movie, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка при получении фильма", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(movie); err != nil {
		http.Error(w, "ошибка в JSON:", http.StatusInternalServerError)
		return
	}
}

func (h *MovieHandler) CreateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request movieRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Неправильный JSON", http.StatusBadRequest)
		return
	}

	movie, err := h.service.CreateMovie(
		ctx,
		request.Title,
		request.ReleaseYear,
		request.Description,
		request.Duration,
		request.PosterURL,
		request.Director,
		request.Producer,
		request.VideoID,
	)

	if err != nil {
		http.Error(w, "Ошибка в создании фильма", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(movie); err != nil {
		http.Error(w, "ошибка в JSON", http.StatusInternalServerError)
		return
	}

}

func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {

		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteMovie(ctx, id)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка во время удаления фильма", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request movieRequest

	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Неправильный JSON", http.StatusBadRequest)
		return
	}

	req := Movie{
		Title:       request.Title,
		ReleaseYear: request.ReleaseYear,
		Description: request.Description,
		Duration:    request.Duration,
		PosterURL:   request.PosterURL,
		Director:    request.Director,
		Producer:    request.Producer,
		VideoID:     request.VideoID,
	}

	movie, err := h.service.UpdateMovie(ctx, req, id)
	if err != nil {
		if errors.Is(err, ErrMovieNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка во время изменении фильма", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(movie); err != nil {
		http.Error(w, "Ошибка в самом JSON", http.StatusInternalServerError)
		return
	}
}
