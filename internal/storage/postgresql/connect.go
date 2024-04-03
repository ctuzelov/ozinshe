package storage

import (
	"database/sql"
	"fmt"
	"ozinshe/internal/config"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cnf *config.Config) (*Storage, error) {
	const op = "storage.New"

	db, err := sql.Open("postgres", "postgresql://"+cnf.DB.User+":"+cnf.DB.Password+"@"+cnf.Host+":"+cnf.DB.Port+"/"+cnf.DB.Dbname+"?sslmode="+cnf.DB.Sslmode)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(op, "connected successfully")

	return &Storage{db: db}, nil
}
