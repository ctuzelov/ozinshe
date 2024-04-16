package models

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type AgeCategory struct {
	ID     int `json:"id"`
	MinAge int `json:"min_age"`
	MaxAge int `json:"max_age"`
}

type ProjectGenre struct {
	ID          int `json:"id"`
	ProjectID   int `json:"project_id"`
	GenreID     int `json:"genre_id"`
}

type ProjectAgeCategory struct {
	ID            int `json:"id"`
	ProjectID     int `json:"project_id"`
	AgeCategoryID int `json:"age_category_id"`
}

type Keyword struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Cover struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"movie_id"`
	Filename  string `json:"filename"`
}

type Screenshot struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"movie_id"`
	Filename  string `json:"filename"`
}
