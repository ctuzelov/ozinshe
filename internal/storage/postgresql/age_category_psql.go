package storage

import (
	"context"
	"fmt"
	"ozinshe/internal/models"
)

type AgeCategoryStorage struct {
	storage *Postgres
}

func NewAgeCategoryStorage(db *Postgres) *AgeCategoryStorage {
	return &AgeCategoryStorage{storage: db}
}

func (s *AgeCategoryStorage) Insert(ctx context.Context, ageCategories []models.AgeCategory) error {
	const op = "storage.age_category.Insert"

	stmt, err := s.storage.db.Prepare(`
        INSERT INTO age_categories (min_age, max_age) 
        VALUES ($1, $2) 
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

	for _, ac := range ageCategories {
		stmt.ExecContext(ctx, ac.MinAge, ac.MaxAge) // Pass min_age and max_age
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}