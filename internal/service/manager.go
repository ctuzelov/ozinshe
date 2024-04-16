package service

import (
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
)

type User interface {
	Register(user models.User) error
	Login(user models.User) (string, string, error)
	UpdateAllTokens(signedToken string, signedRefreshToken string, user_type string) (string, string, error)
	DeleteTokensByEmail(email string) error
}

type Movie interface {
	Add(movie models.Movie, image_data models.SavePhoto) (int, error)
}

type Series interface {
}

type Genre interface {
	Add(genres []models.Genre) error
}

type AgeCategory interface {
	Add(age_category []models.AgeCategory) error
}

type Keyword interface {
	Add(keywords []models.Keyword) error
}

type Service struct {
	User
	Movie
	Series
	Genre
	Keyword
	AgeCategory
}

func New(storage *psql.Storage) *Service {
	return &Service{
		User:        NewUserService(storage.User),
		Movie:       NewMovieService(storage.Movie),
		Series:      NewSeriesService(storage.Series),
		Genre:       NewGenreService(storage.Genre),
		Keyword:     NewKeywordService(storage.Keyword),
		AgeCategory: NewAgeCategoryService(storage.AgeCategory),
	}
}
