package service

import (
	psql "ozinshe/internal/storage/postgresql"
)

type SeriesService struct {
	Series psql.Series
}

func NewSeriesService(series psql.Series) *SeriesService {
	return &SeriesService{Series: series}
}
