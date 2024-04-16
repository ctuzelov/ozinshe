package storage

import (
	"context"
	"ozinshe/internal/models"
)

type User interface {
	SaveUser(ctx context.Context, user models.User) error
	GetByEmail(ctx context.Context, email string) (models.User, error)
	UpdateTokens(signedToken string, signedRefreshToken string, user_type string) error
	DeleteTokens(ctx context.Context, email string) error
}

type Movie interface {
	Insert(ctx context.Context, movie models.Movie) (int, error)
}

type Series interface {
}

type Genre interface {
	Insert(ctx context.Context, genre []models.Genre) error
}

type AgeCategory interface {
	Insert(ctx context.Context, age_category []models.AgeCategory) error
}

type Keyword interface {
	Insert(ctx context.Context, keywords []models.Keyword) error
}

type Storage struct {
	User
	Movie
	Series
	Genre
	AgeCategory
	Keyword
}

func New(storage *Postgres) *Storage {
	return &Storage{
		User:        NewUserStorage(storage),
		Genre:       NewGenreStorage(storage),
		Movie:       NewMovieStorage(storage),
		Series:      NewSeriesStorage(storage),
		AgeCategory: NewAgeCategoryStorage(storage),
		Keyword:     NewKeywordStorage(storage),
	}
}
