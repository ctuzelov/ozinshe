package models

type Movie struct {
    ID              int    `json:"id"`
    Title           string `json:"title"`
    ReleaseYear     int    `json:"release_year"`
    Description     string `json:"description"`
    Popularity      int    `json:"popularity"`
    YoutubeID       string `json:"youtube_id"`
    Duration        int    `json:"duration"`
    Director        string `json:"director"`
    Producer        string `json:"producer"`
}

type MovieCover struct {
    ID          int    `json:"id"`
    MovieID     int    `json:"movie_id"`
    Name        string `json:"name"`
    Filename    string `json:"filename"`
}

type MovieScreenshot struct {
	ID          int    `json:"id"`
	MovieID     int    `json:"movie_id"`
	Name        string `json:"name"`
	Filename    string `json:"filename"`
}

type MovieAgeCategory struct {
	MovieID      int `json:"movie_id"`
	AgeCategoryID  int `json:"age_category_id"`
}

type MovieGenre struct {
    MovieID  int `json:"movie_id"`
    GenreID  int `json:"genre_id"`
}

type FavoriteMovie struct {
    ID       int `json:"id"`
    UserID   int `json:"user_id"`
    MovieID  int `json:"movie_id"`
}