package movie_genres

import (
	"context"
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
