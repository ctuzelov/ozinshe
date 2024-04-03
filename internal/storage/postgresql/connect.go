package storage

import (
	"database/sql"
	"fmt"
	"ozinshe/internal/config"

	_ "github.com/lib/pq"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(cnf *config.Config) (*Postgres, error) {
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

	return &Postgres{db: db}, nil
}
