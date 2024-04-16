package models

type Series struct {
	ID            int           `json:"id"`
	Title         string        `json:"title"`
	ReleaseYear   int           `json:"release_year"`
	Description   string        `json:"description"`
	Popularity    int           `json:"popularity"`
	Duration      int           `json:"duration"`
	Director      string        `json:"director"`
	Producer      string        `json:"producer"`
	Genres        []Genre       `json:"genres"`
	Cover         Cover         `json:"cover"`
	Screenshots   []Screenshot  `json:"screenshots"`
	AgeCategories []AgeCategory `json:"age_categories"`
}

type Season struct {
	ID           int       `json:"id"`
	SeriesID     int       `json:"series_id"`
	SeasonNumber int       `json:"season_number"`
	Episodes     []Episode `json:"episodes"`
}

type Episode struct {
	ID            int    `json:"id"`
	SeasonID      int    `json:"season_id"`
	EpisodeNumber int    `json:"episode_number"`
	YoutubeID     string `json:"youtube_id"`
}

type FavoriteSeries struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id"`
	MovieID int `json:"movie_id"`
}
