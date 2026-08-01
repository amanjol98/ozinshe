package movie_screenshots

import (
	"context"
)

type MovieScreenshotsService struct {
	repo *MovieScreenshotsRepository
}

func NewMovieScreenshotsService(repo *MovieScreenshotsRepository) *MovieScreenshotsService {
	return &MovieScreenshotsService{repo: repo}
}

func (s *MovieScreenshotsService) AddScreenshotsToMovie(ctx context.Context, movieID int, imageURL string) error {
	return s.repo.AddScreenshotsToMovie(ctx, movieID, imageURL)
}

func (s *MovieScreenshotsService) GetScreenshots(ctx context.Context, movieID int) ([]MovieScreenshot, error) {
	return s.repo.GetScreenshots(ctx, movieID)
}

func (s *MovieScreenshotsService) DeleteScreenshot(ctx context.Context, screenshotID int) error {
	return s.repo.DeleteScreenshot(ctx, screenshotID)
}
