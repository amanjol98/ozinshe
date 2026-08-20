package favorites

type Favorite struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	MovieID int `json:"movie_id"`
}
