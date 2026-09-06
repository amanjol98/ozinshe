package services

import (
	"context"
	"errors"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"

	"strings"
)

type MovieService struct {
	repo *repositories.MovieRepository

	genreService    *MovieGenreService
	categoryService *MovieCategoryService

	seasonService           *SeasonService
	episodeService          *EpisodeService
	movieScreenshotsService *MovieScreenshotsService
}

func NewMovieService(
	repo *repositories.MovieRepository,
	genreService *MovieGenreService,
	categoryService *MovieCategoryService,
	seasonService *SeasonService,
	episodeService *EpisodeService,
	movieScreenshotsService *MovieScreenshotsService,
) *MovieService {
	return &MovieService{
		repo:                    repo,
		genreService:            genreService,
		categoryService:         categoryService,
		seasonService:           seasonService,
		episodeService:          episodeService,
		movieScreenshotsService: movieScreenshotsService,
	}
}

var ErrEmptyTitle = errors.New("Ввели пустое значение")

func (s *MovieService) GetAll(
	ctx context.Context,
	search string,
	categoryID, genreID *int,
	limit, offset int) ([]models.Movie, error) {
	search = strings.TrimSpace(search)
	return s.repo.GetAll(ctx, search, categoryID, genreID, limit, offset)
}

func (s *MovieService) GetHome(ctx context.Context) ([]models.Movie, error) {
	return s.repo.GetHome(ctx)
}

func (s *MovieService) GetByID(ctx context.Context, id int) (models.MovieResponse, error) {
	movie, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.MovieResponse{}, err
	}
	genres, err := s.genreService.GetGenresOfMovie(ctx, id)
	if err != nil {
		return models.MovieResponse{}, err
	}
	categories, err := s.categoryService.GetMovieCategories(ctx, id)
	if err != nil {
		return models.MovieResponse{}, err
	}

	seasons, err := s.seasonService.GetSeasons(ctx, id)
	if err != nil {
		return models.MovieResponse{}, err
	}

	for i := range seasons {
		episodes, err := s.episodeService.GetEpisodes(ctx, seasons[i].ID)
		if err != nil {
			return models.MovieResponse{}, err
		}

		seasons[i].Episodes = episodes
	}

	screenshots, err := s.movieScreenshotsService.GetScreenshots(ctx, id)
	if err != nil {
		return models.MovieResponse{}, err
	}

	response := models.MovieResponse{
		Movie:            movie,
		Genres:           genres,
		Categories:       categories,
		Seasons:          seasons,
		MovieScreenshots: screenshots,
	}
	return response, nil
}

func (s *MovieService) CreateMovie(
	ctx context.Context,
	title string,
	releaseYear int,
	description string,
	duration int,
	posterURL string,
	director string,
	producer string,
	videoID string,
	genreIDs,
	categoryIDs []int,
) (models.CreateMovieModel, error) {
	if strings.TrimSpace(title) == "" {
		return models.CreateMovieModel{}, ErrEmptyTitle
	}

	createdMovie, err := s.repo.CreateMovie(
		ctx,
		title,
		releaseYear,
		description,
		duration,
		posterURL,
		director,
		producer,
		videoID,
	)

	if err != nil {
		return models.CreateMovieModel{}, err
	}

	for _, genreID := range genreIDs {
		err := s.genreService.AddGenreToMovie(ctx, createdMovie.ID, genreID)
		if err != nil {
			return models.CreateMovieModel{}, err
		}
	}

	for _, categoryID := range categoryIDs {
		err := s.categoryService.AddCategoryToMovie(ctx, createdMovie.ID, categoryID)
		if err != nil {
			return models.CreateMovieModel{}, err
		}
	}

	genres, err := s.genreService.GetGenresOfMovie(ctx, createdMovie.ID)
	if err != nil {
		return models.CreateMovieModel{}, err
	}
	categories, err := s.categoryService.GetMovieCategories(ctx, createdMovie.ID)
	if err != nil {
		return models.CreateMovieModel{}, err
	}

	return models.CreateMovieModel{
		Movie:      createdMovie,
		Genres:     genres,
		Categories: categories,
	}, nil

}

func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
	return s.repo.DeleteMovie(ctx, id)
}

func (s *MovieService) UpdateMovie(
	ctx context.Context,
	movie models.Movie,
	genreIDs []int,
	categoryIDs []int,
	id int,
) (models.CreateMovieModel, error) {
	if strings.TrimSpace(movie.Title) == "" {
		return models.CreateMovieModel{}, ErrEmptyTitle
	}

	updatedMovie, err := s.repo.UpdateMovie(ctx, movie, id)
	if err != nil {
		return models.CreateMovieModel{}, err
	}

	err = s.genreService.DeleteAllGenresFromMovie(ctx, id)
	if err != nil {
		return models.CreateMovieModel{}, err
	}

	for _, genreID := range genreIDs {
		err := s.genreService.AddGenreToMovie(ctx, id, genreID)
		if err != nil {
			return models.CreateMovieModel{}, err
		}
	}

	err = s.categoryService.DeleteAllCategoriesFromMovie(ctx, id)
	if err != nil {
		return models.CreateMovieModel{}, err
	}

	for _, categoryID := range categoryIDs {
		err := s.categoryService.AddCategoryToMovie(ctx, id, categoryID)
		if err != nil {
			return models.CreateMovieModel{}, err
		}
	}

	genres, err := s.genreService.GetGenresOfMovie(ctx, id)
	if err != nil {
		return models.CreateMovieModel{}, err
	}
	categories, err := s.categoryService.GetMovieCategories(ctx, id)
	if err != nil {
		return models.CreateMovieModel{}, err
	}

	return models.CreateMovieModel{
		Movie:      updatedMovie,
		Genres:     genres,
		Categories: categories,
	}, nil
}
