package storage

import (
	"context"
	"fmt"
	"ozinshe/internal/models"
)

type MovieStorage struct {
	storage *Postgres
}

func NewMovieStorage(db *Postgres) *MovieStorage {
	return &MovieStorage{storage: db}
}

func (m *MovieStorage) Insert(ctx context.Context, movie models.Movie) (int, error) {
	const op = "storage.movie.Insert"

	tx, err := m.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	// Insert Movie
	var movieID int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO movies (title, release_year, description, popularity, youtube_id, duration, director, producer)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
         RETURNING id`,
		movie.Title, movie.ReleaseYear, movie.Description, movie.Popularity,
		movie.YoutubeID, movie.Duration, movie.Director, movie.Producer,
	).Scan(&movieID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: insert movie: %w", op, err)
	}

	// Insert Genres
	for _, genre := range movie.Genres {
		// 1. Attempt to Insert the Genre
		_, err := tx.ExecContext(ctx,
			`INSERT INTO genres (name) 
			 VALUES ($1)
			 ON CONFLICT (name) DO NOTHING`, // Key change: ON CONFLICT
			genre.Name,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: insert genre (attempt): %w", op, err)
		}

		// 2. Get Genre ID (Inserted or Existing)
		var genreID int
		err = tx.QueryRowContext(ctx,
			`SELECT id FROM genres WHERE name = $1`,
			genre.Name,
		).Scan(&genreID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: get genre id: %w", op, err)
		}

		// 3. Link Genre with Movie
		_, err = tx.ExecContext(ctx,
			`INSERT INTO movie_genres (movie_id, genre_id) 
			 VALUES ($1, $2)`,
			movieID, genreID,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: link genre to movie: %w", op, err)
		}
	}

	// Insert Keywords
	for _, keyword := range movie.Keywords {
		// 1. Attempt to Insert Keyword
		_, err := tx.ExecContext(ctx,
			`INSERT INTO key_words (name) 
		 VALUES ($1)
		 ON CONFLICT (name) DO NOTHING`,
			keyword.Name,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: insert keyword (attempt): %w", op, err)
		}

		// 2. Get Keyword ID
		var keywordID int
		err = tx.QueryRowContext(ctx,
			`SELECT id FROM key_words WHERE name = $1`,
			keyword.Name,
		).Scan(&keywordID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: get keyword id: %w", op, err)
		}

		// 3. Link Keyword with Movie
		_, err = tx.ExecContext(ctx,
			`INSERT INTO movie_key_words (movie_id, key_word_id) 
		 VALUES ($1, $2)`,
			movieID, keywordID,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: link keyword to movie: %w", op, err)
		}
	}

	// Insert Age Categories
	for _, ageCategory := range movie.AgeCategories {
		// 1. Attempt to Insert Age Category
		_, err := tx.ExecContext(ctx,
			`INSERT INTO age_categories (min_age, max_age) 
		 VALUES ($1, $2)
		 ON CONFLICT (min_age, max_age) DO NOTHING`,
			ageCategory.MinAge, ageCategory.MaxAge,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: insert age category (attempt): %w", op, err)
		}

		// 2. Get Age Category ID
		var ageCategoryID int
		err = tx.QueryRowContext(ctx,
			`SELECT id FROM age_categories WHERE min_age = $1 AND max_age = $2`,
			ageCategory.MinAge, ageCategory.MaxAge,
		).Scan(&ageCategoryID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: get age category id: %w", op, err)
		}

		// 3. Link Age Category with Movie
		_, err = tx.ExecContext(ctx,
			`INSERT INTO movie_age_categories (movie_id, age_category_id) 
		 VALUES ($1, $2)`,
			movieID, ageCategoryID,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: link age category to movie: %w", op, err)
		}
	}
	// Insert Screenshots
	for _, screenshot := range movie.Screenshots {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO movie_screenshots (movie_id, filename) 
             VALUES ($1, $2)`,
			movieID, screenshot.Filename,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: insert screenshots: %w", op, err)
		}
	}

	// Insert Cover
	_, err = tx.ExecContext(ctx,
		`INSERT INTO movie_covers (movie_id, filename) 
         VALUES ($1, $2)`,
		movieID, movie.Cover.Filename,
	)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: insert cover: %w", op, err)
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return movieID, nil
}
