package services

import (
	"context"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"
)

type MovieScreenshotsService struct {
	repo *repositories.MovieScreenshotsRepository
}

func NewMovieScreenshotsService(repo *repositories.MovieScreenshotsRepository) *MovieScreenshotsService {
	return &MovieScreenshotsService{repo: repo}
}

func (s *MovieScreenshotsService) AddScreenshotsToMovie(ctx context.Context, movieID int, imageURL string) error {
	return s.repo.AddScreenshotsToMovie(ctx, movieID, imageURL)
}

func (s *MovieScreenshotsService) GetScreenshots(ctx context.Context, movieID int) ([]models.MovieScreenshot, error) {
	return s.repo.GetScreenshots(ctx, movieID)
}

func (s *MovieScreenshotsService) DeleteScreenshot(ctx context.Context, screenshotID int) error {
	return s.repo.DeleteScreenshot(ctx, screenshotID)
}
