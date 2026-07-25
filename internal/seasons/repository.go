package seasons

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type SeasonRepository struct {
	db *pgx.Conn
}

func NewSeasonRepository(db *pgx.Conn) *SeasonRepository {
	return &SeasonRepository{db: db}
}

func (r *SeasonRepository) GetSeasons(ctx context.Context, movieID int) ([]Season, error) {
	sqlQuery := `
	SELECT
	id,
	season_number
	FROM seasons
	JOIN movie_id ON seasons.movie_id=movies.id
	WHERE movies.id=$1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seasons []Season

	for rows.Next() {
		var season Season
		err := rows.Scan(
			&season.ID,
			season.SeasonNumber,
		)

		if err != nil {
			return nil, err
		}

		seasons = append(seasons, season)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seasons, nil
}
