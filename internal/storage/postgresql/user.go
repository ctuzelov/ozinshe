package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ozinshe/internal/models"
	"ozinshe/internal/storage"
	"time"

	"github.com/lib/pq"
)

type UserStorage struct {
	storage *Postgres
}

func NewUserStorage(db *Postgres) *UserStorage {
	return &UserStorage{storage: db}
}

func (u *UserStorage) SaveUser(ctx context.Context, user models.User) (err error) {
	const op = "storage.user.SaveUser"

	user.Name = ""
	user.Number = ""
	user.DateOfBirth = time.Time{}
	user.UserType = ""

	stmt, err := u.storage.db.Prepare("INSERT INTO users (name, email, number, date_of_birth, user_type, password, token, refresh_token) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, user.Name, user.Email, user.Number, user.DateOfBirth, user.UserType, user.Password, user.Token, user.Refresh_Token)

	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return
}

func (u *UserStorage) GetByEmail(ctx context.Context, email string) (models.User, error) {
	const op = "storage.user.GetByEmail"

	stmt, err := u.storage.db.Prepare("SELECT id, name, email, number, date_of_birth, user_type, password FROM users WHERE email = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Number, &user.DateOfBirth, &user.UserType, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *UserStorage) UpdateTokens(signedToken string, signedRefreshToken string, user_type string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	const op = "storage.user.UpdateTokens"

	stmt, err := u.storage.db.Prepare(`UPDATE users SET token = $1, refresh_token = $2 WHERE user_type = $3`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, signedToken, signedRefreshToken, user_type)

	if err := row.Err(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
