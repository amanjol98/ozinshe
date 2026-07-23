package movie

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
