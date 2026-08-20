package services

import (
	"context"
	"ozinshe/internal/repositories"
)

type FavoriteService struct {
	repo *repositories.FavoriteRepository
}

func NewFavoriteService(repo *repositories.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) AddFavoriteMovieToUser(ctx context.Context, userID, movieID int) error {
	return s.repo.AddFavoriteMovieToUser(ctx, userID, movieID)
}
