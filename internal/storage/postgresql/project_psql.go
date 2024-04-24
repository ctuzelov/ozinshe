package storage

import (
	"context"
	"fmt"
	"ozinshe/internal/models"
)

type ProjectStorage struct {
	storage *Postgres
}

func NewProjectStorage(db *Postgres) *ProjectStorage {
	return &ProjectStorage{storage: db}
}

func (p *ProjectStorage) Insert(ctx context.Context, project models.Project) (int, error) {
	const op = "storage.project.Insert"

	stmt, err := p.storage.db.Prepare(`INSERT INTO projects (project_type, project_id) VALUES ($1, $2) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var id int
	err = stmt.QueryRowContext(ctx, project.Project_type, project.Project_id).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: insert project: %w", op, err)
	}

	return id, nil
}

func (p *ProjectStorage) Delete(ctx context.Context, id int) error {
	const op = "storage.project.Delete"

	stmt, err := p.storage.db.Prepare(`DELETE FROM projects WHERE id = $1`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: delete project: %w", op, err)
	}

	return nil
}

func (p *ProjectStorage) GetById(ctx context.Context, id int) (models.Project, error) {
	const op = "storage.project.GetById"

	stmt, err := p.storage.db.Prepare(`SELECT project_type, project_id FROM projects WHERE id = $1`)
	if err != nil {
		return models.Project{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var project models.Project
	err = stmt.QueryRowContext(ctx, id).Scan(&project.Project_type, &project.Project_id)
	if err != nil {
		return models.Project{}, fmt.Errorf("%s: get project: %w", op, err)
	}

	return project, nil
}

func (p *ProjectStorage) InsertToFavorites(ctx context.Context, movieID, userID int) error {
	const op = "storage.project.InsertToFavorites"

	stmt, err := p.storage.db.Prepare(`INSERT INTO favorite_projects (user_id, project_id) VALUES ($1, $2)	ON CONFLICT DO NOTHING`)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, movieID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ProjectStorage) DeleteFromFavorites(ctx context.Context, movieID, userID int) error {
	const op = "storage.project.DeleteFromFavorites"

	stmt, err := p.storage.db.Prepare(`DELETE FROM favorite_projects WHERE user_id = $1 AND project_id = $2`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, movieID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
