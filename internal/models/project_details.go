package models

const Admin = "chingizkhan.tuzelov@gmail.com"

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
	ID        int `json:"id"`
	ProjectID int `json:"project_id"`
	GenreID   int `json:"genre_id"`
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
	ProjectID int    `json:"project_id"`
	Filename  string `json:"filename"`
}

type Screenshot struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Filename  string `json:"filename"`
}

type Project struct {
	Id           int    `json:"id"`
	Project_type string `json:"project_type"`
	Project_id   int    `json:"project_id"`
	Movies       []Movie
	Series       []Series
}

type FilterParams struct {
	Project_type    string   `form:"project_type"`
	Genres          []string `form:"genres"`
	Age             []string `form:"age"`
	Title           string   `form:"title"`
	YearStart       int      `form:"year_start"`
	YearEnd         int      `form:"year_end"`
	YearOrder       string   `form:"year_order"`
	PopularityOrder string   `form:"popularity_order"`
}
