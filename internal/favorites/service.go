package favorites

import "context"

type FavoriteService struct {
	repo *FavoriteRepository
}

func NewFavoriteService(repo *FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) AddFavoriteMovieToUser(ctx context.Context, userID, movieID int) error {
	return s.repo.AddFavoriteMovieToUser(ctx, userID, movieID)
}
