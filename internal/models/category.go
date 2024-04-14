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

type Keyword struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
