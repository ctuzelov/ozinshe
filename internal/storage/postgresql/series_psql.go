package storage

type SeriesStorage struct {
	storage *Postgres
}

func NewSeriesStorage(db *Postgres) *SeriesStorage {
	return &SeriesStorage{storage: db}
}
