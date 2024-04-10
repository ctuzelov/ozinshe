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

type Service struct {
	User
}

func New(storage *psql.Storage) *Service {
	return &Service{
		User: NewUserService(storage.User),
		// TODO: implement user service init
	}
}
