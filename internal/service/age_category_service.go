package service

import (
	"context"
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
	"time"
)

type AgeCategoryService struct {
	AgeCategory psql.AgeCategory
}

func NewAgeCategoryService(ageCategory psql.AgeCategory) *AgeCategoryService {
	return &AgeCategoryService{AgeCategory: ageCategory}
}

func (a *AgeCategoryService) Add(age_category []models.AgeCategory) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return a.AgeCategory.Insert(ctx, age_category)
}
