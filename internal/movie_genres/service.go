package movie_genres

import (
	"context"
	"ozinshe/internal/genres"
)

type MovieGenreService struct {
	repo *MovieGenreRepository
}

func NewMovieGenreService(repo *MovieGenreRepository) *MovieGenreService {
	return &MovieGenreService{repo: repo}
}

func (s *MovieGenreService) AddGenreToMovie(ctx context.Context, movieID, genreID int) error {
	return s.repo.AddGenreToMovie(ctx, movieID, genreID)
}

func (s *MovieGenreService) GetGenresOfMovie(ctx context.Context, movieID int) ([]genres.Genre, error) {
	return s.repo.GetGenresOfMovie(ctx, movieID)
}

func (s *MovieGenreService) DeleteGenreFromMovie(ctx context.Context, movieID, genreID int) error {
	return s.repo.DeleteGenreFromMovie(ctx, movieID, genreID)
}
