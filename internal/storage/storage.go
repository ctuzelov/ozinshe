package storage

import "fmt"

var (
	ErrUserNotFound = fmt.Errorf("user not found")
	ErrUserExists   = fmt.Errorf("user already exists")
	ErrMovieExists  = fmt.Errorf("movie already exists")
	ErrSeriesExists = fmt.Errorf("series already exists")
)
