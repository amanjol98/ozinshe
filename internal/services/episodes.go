package services

import (
	"context"
	"ozinshe/internal/models"
	"ozinshe/internal/repositories"
)

type EpisodeService struct {
	repo *repositories.EpisodeRepository
}

func NewEpisodeService(repo *repositories.EpisodeRepository) *EpisodeService {
	return &EpisodeService{repo: repo}
}

func (s *EpisodeService) AddEpisodeToSeason(ctx context.Context, seasonID, episodeNumber int, video_id string) error {
	return s.repo.AddEpisodeToSeason(ctx, seasonID, episodeNumber, video_id)
}

func (s *EpisodeService) GetEpisodes(ctx context.Context, seasonID int) ([]models.Episode, error) {
	return s.repo.GetEpisodes(ctx, seasonID)
}

func (s *EpisodeService) DeleteEpisode(ctx context.Context, episodeID int) error {
	return s.repo.DeleteEpisode(ctx, episodeID)
}
