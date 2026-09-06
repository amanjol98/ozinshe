package services

import (
	"context"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"
)

type MovieGenreService struct {
	repo *repositories.MovieGenreRepository
}

func NewMovieGenreService(repo *repositories.MovieGenreRepository) *MovieGenreService {
	return &MovieGenreService{repo: repo}
}

func (s *MovieGenreService) AddGenreToMovie(ctx context.Context, movieID, genreID int) error {
	return s.repo.AddGenreToMovie(ctx, movieID, genreID)
}

func (s *MovieGenreService) GetGenresOfMovie(ctx context.Context, movieID int) ([]models.Genre, error) {
	return s.repo.GetGenresOfMovie(ctx, movieID)
}

func (s *MovieGenreService) DeleteGenreFromMovie(ctx context.Context, movieID, genreID int) error {
	return s.repo.DeleteGenreFromMovie(ctx, movieID, genreID)
}

func (s *MovieGenreService) DeleteAllGenresFromMovie(ctx context.Context, movieID int) error {
	return s.repo.DeleteAllGenresFromMovie(ctx, movieID)
}
