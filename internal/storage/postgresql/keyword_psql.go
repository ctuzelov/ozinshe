package storage

import (
	"context"
	"fmt"
	"ozinshe/internal/models"
)

type KeywordStorage struct {
	storage *Postgres
}

func NewKeywordStorage(db *Postgres) *KeywordStorage {
	return &KeywordStorage{storage: db}
}

func (s *KeywordStorage) Insert(ctx context.Context, keywords []models.Keyword) error {
	const op = "storage.keyword.Insert"

	stmt, err := s.storage.db.Prepare(`
		INSERT INTO key_words (name) 
		VALUES ($1) 
		ON CONFLICT DO NOTHING
	`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	tx, err := s.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback()

	for _, k := range keywords {
		stmt.ExecContext(ctx, k.Name)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}
