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
