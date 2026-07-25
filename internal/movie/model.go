package movie

import (
	"ozinshe/internal/categories"
	"ozinshe/internal/genres"
)

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"release_yesr"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	PosterURL   string `json:"poster_url"`
	Director    string `json:"director"`
	Producer    string `json:"producer"`
	VideoID     string `json:"video_id"`
}

type MovieResponse struct {
	Movie      Movie                 `json:"movie"`
	Genres     []genres.Genre        `json:"genres"`
	Categories []categories.Category `json:"categories"`
}
