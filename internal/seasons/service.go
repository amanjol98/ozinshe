package seasons

import "context"

type SeasonService struct {
	repo *SeasonRepository
}

func NewSeasonService(repo *SeasonRepository) *SeasonService {
	return &SeasonService{repo: repo}
}

func (s *SeasonService) GetSeasons(ctx context.Context, movieID int) ([]Season, error) {
	return s.repo.GetSeasons(ctx, movieID)
}

func (s *SeasonService) AddSeasonToMovie(ctx context.Context, movieID, seasonNumber int) error {
	return s.repo.AddSeasonToMovie(ctx, movieID, seasonNumber)
}

func (s *SeasonService) DeleteSeasonFromMovie(ctx context.Context, seasonID int) error {
	return s.repo.DeleteSeasonFromMovie(ctx, seasonID)
}
