package repositories

import (
	"context"
	"errors"
	"ozinshe/internal/models"

	"github.com/jackc/pgx/v5"
)

type EpisodeRepository struct {
	db *pgx.Conn
}

func NewEpisodeRepository(db *pgx.Conn) *EpisodeRepository {
	return &EpisodeRepository{db: db}
}

var ErrEpisodeNotFound = errors.New("Нет эпизода с таким ID")

func (r *EpisodeRepository) AddEpisodeToSeason(ctx context.Context, seasonID, episodeNumber int, VideoID string) error {
	sqlQuery := `
	INSERT INTO episodes (
	season_id,
	episode_number,
	video_id
	)
	VALUES($1,$2,$3)
	`

	_, err := r.db.Exec(ctx, sqlQuery, seasonID, episodeNumber, VideoID)
	return err
}

func (r *EpisodeRepository) GetEpisodes(ctx context.Context, seasonID int) ([]models.Episode, error) {
	sqlQuery := `
	SELECT
	id,
	episode_number,
	video_id
	FROM episodes
	WHERE season_id=$1;
	`

	rows, err := r.db.Query(ctx, sqlQuery, seasonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var episodes []models.Episode

	for rows.Next() {
		var episode models.Episode
		err := rows.Scan(
			&episode.ID,
			&episode.EpisodeNumber,
			&episode.VideoID,
		)
		if err != nil {
			return nil, err
		}

		episodes = append(episodes, episode)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return episodes, nil
}

func (r *EpisodeRepository) DeleteEpisode(ctx context.Context, episodeID int) error {
	sqlQuery := `
	DELETE FROM episodes
	WHERE id=$1;
	`

	result, err := r.db.Exec(ctx, sqlQuery, episodeID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrEpisodeNotFound
	}
	return err
}
