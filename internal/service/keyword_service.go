package service

import (
	"context"
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
	"time"
)

type KeywordService struct {
	Keyword psql.Keyword
}

func NewKeywordService(keyword psql.Keyword) *KeywordService {
	return &KeywordService{Keyword: keyword}
}

func (k *KeywordService) Add(keywords []models.Keyword) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return k.Keyword.Insert(ctx, keywords)
}
