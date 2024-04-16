package service

import (
	"context"
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
	"time"
)

type GenreService struct {
	Genre psql.Genre
}

func NewGenreService(genre psql.Genre) *GenreService {
	return &GenreService{Genre: genre}
}

func (g *GenreService) Add(genre []models.Genre) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return g.Genre.Insert(ctx, genre)
}
