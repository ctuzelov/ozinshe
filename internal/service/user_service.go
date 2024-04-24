package service

import (
	"context"
	"fmt"
	"ozinshe/internal/helper"
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
	"ozinshe/util"
	"strconv"
	"time"
)

var (
	ErrWrongPassword = fmt.Errorf("wrong password")
)

type UserService struct {
	Storage psql.User
}

func NewUserService(storage psql.User) *UserService {
	return &UserService{Storage: storage}
}

func (a *UserService) Register(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	const op = "service.user.Register"

	user.Password = helper.HashPassword(user.Password)
	user.Token, user.Refresh_Token, err = util.GenerateAllTokens(user.Email, user.Name, user.UserType, strconv.Itoa(user.ID))

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.Storage.SaveUser(ctx, user)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *UserService) Login(user models.User) (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.user.Login"
	var foundUser models.User

	foundUser, err := a.Storage.GetByEmail(ctx, user.Email)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	match, _ := helper.VerifyPassword(user.Password, foundUser.Password)
	if !match {
		return "", "", fmt.Errorf("%s: %w", op, ErrWrongPassword)
	}

	if user.Email == models.Admin{
		foundUser.UserType = "admin"
	}

	foundUser.Token, foundUser.Refresh_Token, err = util.GenerateAllTokens(foundUser.Email, foundUser.Name, foundUser.UserType, strconv.Itoa(foundUser.ID))
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	token, refreshToken, err := a.UpdateAllTokens(foundUser.Token, foundUser.Refresh_Token, foundUser.UserType, strconv.Itoa(foundUser.ID))
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return token, refreshToken, nil
}

func (a *UserService) UpdateAllTokens(signedToken string, signedRefreshToken string, user_type string, id string) (string, string, error) {
	const op = "util.UpdateAllTokens"

	_, err := util.ValidateToken(signedRefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	err = a.Storage.UpdateTokens(signedToken, signedRefreshToken, user_type)
	if err != nil {
		return "", "", fmt.Errorf("%s: %w", op, err)
	}

	return signedToken, signedRefreshToken, nil
}

func (a *UserService) DeleteTokensByEmail(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	const op = "service.user.DeleteTokensByEmail"

	err := a.Storage.DeleteTokens(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *UserService) GetById(id string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.user.GetById"

	user, err := a.Storage.GetByEmail(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (a *UserService) GetAll() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.user.GetAll"

	users, err := a.Storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (a *UserService) UpdatePassword(email string, current_password, new_password string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.user.UpdatePassword"

	user, err := a.Storage.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if user.Password != helper.HashPassword(current_password) {
		return fmt.Errorf("%s: %w", op, ErrWrongPassword)
	}
	user.Password = helper.HashPassword(new_password)
	err = a.Storage.ChangePassword(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *UserService) UpdateProfile(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.user.UpdateProfile"

	user, err := a.Storage.GetByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if user.Password != helper.HashPassword(user.Password) {
		return fmt.Errorf("%s: %w", op, ErrWrongPassword)
	}

	err = a.Storage.ChangeProfileData(ctx, user)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *UserService) Remove(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.user.Remove"

	err := a.Storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
