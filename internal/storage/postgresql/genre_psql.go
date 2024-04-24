package storage

import (
	"context"
	"fmt"
	"ozinshe/internal/models"
)

type GenreStorage struct {
	Storage *Postgres
}

func NewGenreStorage(db *Postgres) *GenreStorage {
	return &GenreStorage{Storage: db}
}

func (s *GenreStorage) Insert(ctx context.Context, genres []models.Genre) error {
	const op = "storage.genre.Insert"

	stmt, err := s.Storage.db.Prepare(`
        INSERT INTO genres (name) 
        VALUES ($1) 
        ON CONFLICT DO NOTHING 
    `)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	tx, err := s.Storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	for _, g := range genres {
		stmt.ExecContext(ctx, g.Name)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *GenreStorage) Delete(ctx context.Context, id int) error {
	const op = "storage.genre.Delete"

	stmt, err := s.Storage.db.Prepare(`DELETE FROM genres WHERE id = $1`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *GenreStorage) GetById(ctx context.Context, id int) (models.Genre, error) {
	const op = "storage.genre.GetById"

	stmt, err := s.Storage.db.Prepare(`
		SELECT id, name
		FROM genres
		WHERE id = $1
	`)
	if err != nil {
		return models.Genre{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var g models.Genre
	err = stmt.QueryRowContext(ctx, id).Scan(&g.ID, &g.Name)
	if err != nil {
		return models.Genre{}, fmt.Errorf("%s: %w", op, err)
	}

	return g, nil
}

func (s *GenreStorage) GetAll(ctx context.Context) ([]models.Genre, error) {
	const op = "storage.genre.GetAll"

	stmt, err := s.Storage.db.Prepare(`
		SELECT id, name
		FROM genres
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		genres = append(genres, g)
	}
	return genres, nil
}

func (s *GenreStorage) GetByName(ctx context.Context, name string) (models.Genre, error) {
	const op = "storage.genre.GetByName"

	stmt, err := s.Storage.db.Prepare(`
		SELECT id, name
		FROM genres
		WHERE name = $1
	`)
	if err != nil {
		return models.Genre{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var g models.Genre
	err = stmt.QueryRowContext(ctx, name).Scan(&g.ID, &g.Name)
	if err != nil {
		return models.Genre{}, fmt.Errorf("%s: %w", op, err)
	}

	return g, nil
}
