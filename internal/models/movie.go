package models

type Movie struct {
	ID            int           `json:"id"`
	Title         string        `json:"title"`
	ReleaseYear   int           `json:"release_year"`
	Description   string        `json:"description"`
	Popularity    int           `json:"popularity"`
	YoutubeID     string        `json:"youtube_id"`
	Duration      int           `json:"duration"`
	Director      string        `json:"director"`
	Producer      string        `json:"producer"`
	Genres        []Genre       `json:"genres"`
	Keywords      []Keyword     `json:"keywords"`
	AgeCategories []AgeCategory `json:"age_categories"`
	Screenshots   []Screenshot  `json:"screenshots"`
	Cover         Cover         `json:"cover"`
}

type FavoriteMovie struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	MovieID int `json:"movie_id"`
}
