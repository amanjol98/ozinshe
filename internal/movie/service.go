package movie

import (
	"context"
	"fmt"
	"strings"
)

type MovieService struct {
	repo *MovieRepository
}

func NewMovieService(repo *MovieRepository) *MovieService {
	return &MovieService{
		repo: repo,
	}
}

func (s *MovieService) GetAll(ctx context.Context) ([]Movie, error) {
	return s.repo.GetAll(ctx)
}

func (s *MovieService) GetByID(ctx context.Context, id int) (Movie, error) {
	return s.repo.GetByID(ctx, id)
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
) (Movie, error) {
	if strings.TrimSpace(title) == "" {
		return Movie{}, fmt.Errorf("Ввели пустое значение")
	}

	movie := Movie{
		Title:       title,
		ReleaseYear: releaseYear,
		Description: description,
		Duration:    duration,
		PosterURL:   posterURL,
		Director:    director,
		Producer:    producer,
		VideoID:     videoID,
	}

	return s.repo.CreateMovie(
		ctx,
		movie.Title,
		movie.ReleaseYear,
		movie.Description,
		movie.Duration,
		movie.PosterURL,
		movie.Director,
		movie.Producer,
		movie.VideoID,
	)

}

func (s *MovieService) DeleteMovie(ctx context.Context, id int) error {
	return s.repo.DeleteMovie(ctx, id)
}

func (s *MovieService) UpdateMovie(ctx context.Context, movie Movie, id int) (Movie, error) {
	if strings.TrimSpace(movie.Title) == "" {
		return Movie{}, fmt.Errorf("Ввели пустое значение")
	}

	return s.repo.UpdateMovie(ctx, movie, id)
}
