package models

type Series struct {
    ID              int    `json:"id"`
    Title           string `json:"title"`
    ReleaseYear     int    `json:"release_year"`
    Description     string `json:"description"`
    Popularity      int    `json:"popularity"`
    Duration        int    `json:"duration"`
    Director        string `json:"director"`
    Producer        string `json:"producer"`
}

type Season struct {
    ID           int `json:"id"`
    SeriesID     int `json:"series_id"`
    SeasonNumber int `json:"season_number"`
}

type Episode struct {
    ID           int    `json:"id"`
    SeasonID     int    `json:"season_id"`
    EpisodeNumber int   `json:"episode_number"` 
    YoutubeID    string `json:"youtube_id"`
}

type SeriesGenre struct {
    MovieID  int `json:"movie_id"`
    GenreID  int `json:"genre_id"`
}

type SeriesAgeCategory struct {
	MovieID      int `json:"movie_id"`
	AgeCategoryID  int `json:"age_category_id"`
}

type SeriesCover struct {
	ID          int    `json:"id"`
	SeriesID     int    `json:"series_id"`
	Name        string `json:"name"`
	Filename    string `json:"filename"`
}

type SeriesScreenshot struct {
	ID          int    `json:"id"`
	SeriesID     int    `json:"series_id"`
	Name        string `json:"name"`
	Filename    string `json:"filename"`
}

type FavoriteSeries struct {
    ID       int `json:"id"`
    UserID   int `json:"user_id"`
    MovieID  int `json:"movie_id"`
}