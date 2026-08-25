package repositories

import (
	"context"
	"ozinshe/internal/models"

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

func (r *FavoriteRepository) GetFavoriteMovies(ctx context.Context, userID int) ([]models.FavoriteMovieResponse, error) {
	sqlQuery := `
		SELECT
			m.id,
			m.title,
			m.release_year,
			m.poster_url,
			COALESCE(
				ARRAY_AGG(c.name ORDER BY c.name)
				FILTER (WHERE c.id IS NOT NULL),
				'{}'
			) AS categories
		FROM favorites f
		JOIN movies m ON m.id = f.movie_id
		LEFT JOIN movie_categories mc ON mc.movie_id=m.id
		LEFT JOIN categories c ON c.id=mc.category_id
		WHERE f.user_id=&1
		GROUP BY
		m.id,
		m.title,
		m.release_year,
		m.poster_url,
		f.id
		ORDER BY f.id DESC;
	`

	rows, err := r.db.Query(ctx, sqlQuery, userID)
	if err != nil {
		return nil, err
	}

	var favorites []models.FavoriteMovieResponse

	for rows.Next() {
		var movie models.FavoriteMovieResponse

		err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.PosterURL,
			&movie.Categories,
		)

		if err != nil {
			return nil, err
		}

		favorites = append(favorites, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return favorites, nil
}
