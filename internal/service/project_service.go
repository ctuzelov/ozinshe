package service

import (
	"context"
	"fmt"
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
	"time"
)

type ProjectService struct {
	storage psql.Project
}

func NewProjectService(storage psql.Project) *ProjectService {
	return &ProjectService{storage: storage}
}

func (p *ProjectService) Add(project models.Project) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.project.Add"

	id, err := p.storage.Insert(ctx, project)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *ProjectService) Remove(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.project.Remove"

	err := p.storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ProjectService) GetById(id int) (models.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.project.GetById"

	project, err := p.storage.GetById(ctx, id)
	if err != nil {
		return models.Project{}, fmt.Errorf("%s: %w", op, err)
	}

	return project, nil
}

func (p *ProjectService) AddToFavorites(movieID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.project.AddToFavorites"

	err := p.storage.InsertToFavorites(ctx, movieID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *ProjectService) RemoveFromFavorites(movieID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.project.RemoveFromFavorites"

	err := p.storage.DeleteFromFavorites(ctx, movieID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
