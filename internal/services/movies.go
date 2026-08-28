package services

import (
	"context"
	"fmt"
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

func (s *MovieService) GetAll(ctx context.Context, search string, limit, offset int) ([]models.Movie, error) {
	search = strings.TrimSpace(search)
	return s.repo.GetAll(ctx, search, limit, offset)
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
		return models.CreateMovieModel{}, fmt.Errorf("Ввели пустое значение")
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

func (s *MovieService) UpdateMovie(ctx context.Context, movie models.Movie, id int) (models.Movie, error) {
	if strings.TrimSpace(movie.Title) == "" {
		return models.Movie{}, fmt.Errorf("Ввели пустое значение")
	}

	return s.repo.UpdateMovie(ctx, movie, id)
}
