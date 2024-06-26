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

func (u *UserStorage) DeleteTokens(ctx context.Context, email string) error {
	const op = "storage.user.GetByToken"

	stmt, err := u.storage.db.Prepare("UPDATE users SET token = NULL, refresh_token = NULL WHERE email = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserStorage) GetAll(ctx context.Context) ([]models.User, error) {
	const op = "storage.user.GetAll"

	stmt, err := u.storage.db.Prepare("SELECT id, name, email, number, date_of_birth, user_type, password FROM users")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Number, &user.DateOfBirth, &user.UserType, &user.Password); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (u *UserStorage) GetById(ctx context.Context, id int) (models.User, error) {
	const op = "storage.user.GetById"

	stmt, err := u.storage.db.Prepare("SELECT id, name, email, number, date_of_birth, user_type, password FROM users WHERE id = $1")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Number, &user.DateOfBirth, &user.UserType, &user.Password); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *UserStorage) ChangePassword(ctx context.Context, user models.User) error {
	const op = "storage.user.ChangePassword"

	stmt, err := u.storage.db.Prepare("UPDATE users SET password = $1 WHERE email = $2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, user.Password, user.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserStorage) ChangeProfileData(ctx context.Context, user models.User) error {
	const op = "storage.user.ChangeProfileData"

	stmt, err := u.storage.db.Prepare("UPDATE users SET name = $1, number = $2, date_of_birth = $3 WHERE email = $4")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, user.Name, user.Number, user.DateOfBirth, user.Email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (u *UserStorage) Delete(ctx context.Context, id int) error {
	const op = "storage.user.Delete"

	stmt, err := u.storage.db.Prepare("DELETE FROM users WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
