package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"
	"ozinshe/internal/services"
	"strconv"
)

type MovieHandler struct {
	service *services.MovieService
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
	GenreIDs    []int  `json:"genre_ids"`
	CategoryIDs []int  `json:"category_ids"`
}

func NewMovieHandler(service *services.MovieService) *MovieHandler {
	return &MovieHandler{
		service: service,
	}
}

// GetAll godoc
// @Summary Получить список фильмов
// @Description Возвращает список фильмов с поиском, фильтрацией по категории и жанру, а также пагинацией.
// @Tags Movies
// @Produce json
// @Param search query string false "Поиск по названию"
// @Param category_id query int false "ID категории"
// @Param genre_id query int false "ID жанра"
// @Param page query int false "Номер страницы"
// @Param limit query int false "Количество фильмов на странице"
// @Success 200 {array} models.Movie
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies [get]
func (h *MovieHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var err error

	search := r.URL.Query().Get("search")

	var categoryID *int
	var genreID *int

	categoryIDStr := r.URL.Query().Get("category_id")

	if categoryIDStr != "" {
		id, err := strconv.Atoi(categoryIDStr)
		if err != nil || id <= 0 {
			http.Error(w, "Неверный category_id", http.StatusBadRequest)
			return
		}

		categoryID = &id
	}

	genreIDStr := r.URL.Query().Get("genre_id")

	if genreIDStr != "" {
		id, err := strconv.Atoi(genreIDStr)
		if err != nil || id <= 0 {
			http.Error(w, "Неверный genre_id", http.StatusBadRequest)
			return
		}

		genreID = &id
	}

	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")

	limit := 10
	page := 1

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "неверный limit", http.StatusBadRequest)
			return
		}
		if limit < 1 {
			http.Error(w, "limit должен быть больше 0", http.StatusBadRequest)
			return
		}

		if limit > 100 {
			http.Error(w, "максимальный limit 100", http.StatusBadRequest)
			return
		}
	}

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			http.Error(w, "неверный page", http.StatusBadRequest)
			return
		}
		if page < 1 {
			http.Error(w, "page должен быть больше 0", http.StatusBadRequest)
			return
		}
	}

	offset := (page - 1) * limit

	movies, err := h.service.GetAll(ctx, search, categoryID, genreID, limit, offset)
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

// GetHome godoc
// @Summary На главную
// @Description Возвращает список фильмов с главной страницы.
// @Tags Movies
// @Produce json
// @Success 200 {array} models.Movie
// @Failure 500 {object} map[string]string
// @Router /home [get]
func (h *MovieHandler) GetHome(w http.ResponseWriter, r *http.Request) {
	movies, err := h.service.GetHome(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(movies); err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

}

// GetByID godoc
// @Summary Получить фильм по ID
// @Description Возвращает выбранный по ID фильм.
// @Tags Movies
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} models.Movie
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies/{id} [get]
func (h *MovieHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "ID должен быть больше 0", http.StatusBadRequest)
		return
	}

	movie, err := h.service.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrMovieNotFound) {
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

// CreateMovie godoc
// @Summary Создать фильм
// @Description Создает новый фильм. Доступно только администраторам.
// @Tags Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body movieRequest true "Данные для создания"
// @Success 201 {object} models.CreateMovieModel
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies [post]
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
		request.GenreIDs,
		request.CategoryIDs,
	)

	if err != nil {
		if errors.Is(err, services.ErrEmptyTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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

// DeleteMovie godoc
// @Summary Удалить фильм
// @Description Удаляет выбранный фильм. Доступно только администраторам.
// @Tags Movies
// @Security BearerAuth
// @Param id path int true "ID фильма"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies/{id} [delete]
func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idString := r.PathValue("id")

	id, err := strconv.Atoi(idString)
	if err != nil {

		http.Error(w, "неверный ID", http.StatusBadRequest)
		return
	}

	if id <= 0 {
		http.Error(w, "ID должен быть больше 0", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteMovie(ctx, id)
	if err != nil {
		if errors.Is(err, repositories.ErrMovieNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка во время удаления фильма", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateMovie godoc
// @Summary Обновить фильм
// @Description Обновляет выбранный фильм. Доступно только администраторам.
// @Tags Movies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID фильма"
// @Param request body movieRequest true "Данные для обновления"
// @Success 200 {object} models.CreateMovieModel
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /movies/{id} [patch]
func (h *MovieHandler) UpdateMovie(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request movieRequest

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

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {

		http.Error(w, "Неправильный JSON", http.StatusBadRequest)
		return
	}

	req := models.Movie{
		Title:       request.Title,
		ReleaseYear: request.ReleaseYear,
		Description: request.Description,
		Duration:    request.Duration,
		PosterURL:   request.PosterURL,
		Director:    request.Director,
		Producer:    request.Producer,
		VideoID:     request.VideoID,
	}

	movie, err := h.service.UpdateMovie(ctx, req, request.GenreIDs, request.CategoryIDs, id)
	if err != nil {
		if errors.Is(err, repositories.ErrMovieNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, services.ErrEmptyTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)
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
