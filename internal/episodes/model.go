package episodes

type Episode struct {
	ID            int    `json:"id"`
	EpisodeNumber int    `json:"episode_number"`
	VideoID       string `json:"video_id"`
}
