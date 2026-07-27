package seasons

import "ozinshe/internal/episodes"

type Season struct {
	ID           int                `json:"id"`
	SeasonNumber int                `json:"season_number"`
	Episodes     []episodes.Episode `json:"episodes"`
}
