package movie_genres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type MovieGenreRepository struct {
	db *pgx.Conn
}

func NewMovieGenreRepository(db *pgx.Conn) *MovieGenreRepository {
	return &MovieGenreRepository{db: db}
}

func (r *MovieGenreRepository) AddGenreToMovie(ctx context.Context, movieID, genreID int) error {
	sqlQuery := `
	INSERT INTO movie_genres (
	movie_id,
	genre_id
	)
	VALUES ($1, $2);
	`
	_, err := r.db.Exec(ctx, sqlQuery, movieID, genreID)

	return err
}
