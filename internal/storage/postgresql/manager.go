package storage

import (
	"context"
	"ozinshe/internal/models"
)

type User interface {
	SaveUser(ctx context.Context, user models.User) error
	GetByEmail(ctx context.Context, email string) (models.User, error)
	UpdateTokens(signedToken string, signedRefreshToken string, user_type string) error
}

type Storage struct {
	User
}

func New(storage *Postgres) *Storage {
	return &Storage{
		User: NewUserStorage(storage),
	}
}
