package favorites

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type FavoriteRepository struct {
	db *pgx.Conn
}

func NewFavoriteRepository(db *pgx.Conn) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) AddFavoriteMovieToUser(ctx context.Context, userID, movieID int) error {
	sqlQuery := `
	INSERT INTO favorites (user_id, movie_id)
	VALUES($1,$2)
	ON CONFLICT (user_id, movie_id) DO NOTHING;
	`

	_, err := r.db.Exec(ctx, sqlQuery, userID, movieID)
	return err
}
