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

func (s *AgeCategoryStorage) Delete(ctx context.Context, id int) error {
	const op = "storage.age_category.Delete"

	stmt, err := s.storage.db.Prepare(`DELETE FROM age_categories WHERE id = $1`)
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

func (s *AgeCategoryStorage) GetById(ctx context.Context, id int) (models.AgeCategory, error) {
	const op = "storage.age_category.GetById"

	stmt, err := s.storage.db.Prepare(`
		SELECT id, min_age, max_age
		FROM age_categories
		WHERE id = $1

	`)
	if err != nil {
		return models.AgeCategory{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var ac models.AgeCategory
	err = stmt.QueryRowContext(ctx, id).Scan(&ac.ID, &ac.MinAge, &ac.MaxAge)
	if err != nil {
		return models.AgeCategory{}, fmt.Errorf("%s: %w", op, err)
	}
	return ac, nil
}

func (s *AgeCategoryStorage) GetAll(ctx context.Context) ([]models.AgeCategory, error) {
	const op = "storage.age_category.GetAll"

	stmt, err := s.storage.db.Prepare(`
		SELECT id, min_age, max_age
		FROM age_categories
		ORDER BY id
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

	var ageCategories []models.AgeCategory
	for rows.Next() {
		var ac models.AgeCategory
		err = rows.Scan(&ac.ID, &ac.MinAge, &ac.MaxAge)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		ageCategories = append(ageCategories, ac)
	}

	return ageCategories, nil
}
