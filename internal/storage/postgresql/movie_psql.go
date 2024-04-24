package storage

import (
	"context"
	"errors"
	"fmt"
	"ozinshe/internal/models"
	"ozinshe/internal/storage"
	"strconv"
	"strings"

	"github.com/lib/pq"
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
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrMovieExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
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
			`INSERT INTO age_categories (min_age, max_age, range) 
		 VALUES ($1, $2, $3)
		 ON CONFLICT (range) DO NOTHING`,
			ageCategory.MinAge, ageCategory.MaxAge, strconv.Itoa(ageCategory.MinAge)+"-"+strconv.Itoa(ageCategory.MaxAge),
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

func (m *MovieStorage) Delete(ctx context.Context, id int) error {
	const op = "storage.movie.Delete"

	stmt, err := m.storage.db.Prepare("DELETE FROM movies WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: movie prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: delete movie: %w", op, err)
	}

	return nil
}

func (m *MovieStorage) Update(ctx context.Context, movie models.Movie) error {
	const op = "storage.movie.Update"

	tx, err := m.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	// Update Movie
	if movie.Title != "" || movie.ReleaseYear != 0 || movie.Description != "" ||
		movie.Popularity != 0 || movie.YoutubeID != "" || movie.Duration != 0 ||
		movie.Director != "" || movie.Producer != "" {
		_, err := tx.ExecContext(ctx,
			`UPDATE movies 
			SET title = COALESCE($1, title), release_year = COALESCE($2, release_year), 
			description = COALESCE($3, description), popularity = COALESCE($4, popularity), 
			youtube_id = COALESCE($5, youtube_id), duration = COALESCE($6, duration), 
			director = COALESCE($7, director), producer = COALESCE($8, producer)
			WHERE id = $9`,
			movie.Title, movie.ReleaseYear, movie.Description, movie.Popularity,
			movie.YoutubeID, movie.Duration, movie.Director, movie.Producer,
			movie.ID,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: update movie: %w", op, err)
		}
	}

	// Update Age Categories
	if len(movie.AgeCategories) > 0 {
		// Delete existing age categories
		_, err := tx.ExecContext(ctx, "DELETE FROM movie_age_categories WHERE movie_id = $1", movie.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: delete existing age categories: %w", op, err)
		}

		// Insert new age categories
		for _, ageCategory := range movie.AgeCategories {
			// Insert Age Category
			_, err := tx.ExecContext(ctx,
				`INSERT INTO age_categories (min_age, max_age, range) 
				VALUES ($1, $2, $3)
				ON CONFLICT (range) DO NOTHING`,
				ageCategory.MinAge, ageCategory.MaxAge, strconv.Itoa(ageCategory.MinAge)+"-"+strconv.Itoa(ageCategory.MaxAge),
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: insert age category: %w", op, err)
			}

			// Get Age Category ID
			var ageCategoryID int
			err = tx.QueryRowContext(ctx,
				`SELECT id FROM age_categories WHERE min_age = $1 AND max_age = $2`,
				ageCategory.MinAge, ageCategory.MaxAge,
			).Scan(&ageCategoryID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: get age category id: %w", op, err)
			}

			// Link Age Category with Movie
			_, err = tx.ExecContext(ctx,
				`INSERT INTO movie_age_categories (movie_id, age_category_id) 
				VALUES ($1, $2)`,
				movie.ID, ageCategoryID,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: link age category to movie: %w", op, err)
			}
		}
	}

	// Update Genres
	if len(movie.Genres) > 0 {
		// Delete existing genres
		_, err := tx.ExecContext(ctx, "DELETE FROM movie_genres WHERE movie_id = $1", movie.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: delete existing genres: %w", op, err)
		}

		// Insert new genres
		for _, genre := range movie.Genres {
			// Insert Genre
			_, err := tx.ExecContext(ctx,
				`INSERT INTO genres (name) 
				VALUES ($1)
				ON CONFLICT (name) DO NOTHING`,
				genre.Name,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: insert genre: %w", op, err)
			}

			// Get Genre ID
			var genreID int
			err = tx.QueryRowContext(ctx,
				`SELECT id FROM genres WHERE name = $1`,
				genre.Name,
			).Scan(&genreID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: get genre id: %w", op, err)
			}

			// Link Genre with Movie
			_, err = tx.ExecContext(ctx,
				`INSERT INTO movie_genres (movie_id, genre_id) 
				VALUES ($1, $2)`,
				movie.ID, genreID,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: link genre to movie: %w", op, err)
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (m *MovieStorage) UpdateCover(ctx context.Context, movieID int, cover models.Cover) error {
	const op = "storage.movie.UpdateCover"

	tx, err := m.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// 1. Delete existing cover (if any)
	_, err = tx.ExecContext(ctx, `DELETE FROM movie_covers WHERE movie_id = $1`, movieID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: delete error: %w", op, err)
	}

	// 2. Insert the new cover
	_, err = tx.ExecContext(ctx, `
        INSERT INTO movie_covers (movie_id, filename) VALUES ($1, $2)`, movieID, cover.Filename)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: insert error: %w", op, err)
	}

	return tx.Commit()
}

func (m *MovieStorage) UpdateScreenshots(ctx context.Context, movieID int, screenshots []models.Screenshot) error {
	const op = "storage.movie.UpdateScreenshots"

	tx, err := m.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// 1. Delete existing screenshots
	_, err = tx.ExecContext(ctx, `DELETE FROM movie_screenshots WHERE movie_id = $1`, movieID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: delete error: %w", op, err)
	}

	// 2. Insert new screenshots
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO movie_screenshots (movie_id, filename) VALUES ($1, $2)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare error: %w", op, err)
	}
	defer stmt.Close()

	for _, screenshot := range screenshots {
		_, err = stmt.ExecContext(ctx, movieID, screenshot.Filename)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: insert error: %w", op, err)
		}
	}

	return tx.Commit()
}

func (m *MovieStorage) GetAll(ctx context.Context) ([]models.Movie, error) {
	const op = "storage.movie.GetAll"

	movies, err := m.FetchMovieData(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: Fetch movies: %w", op, err)
	}

	for i := range movies {
		err := m.FetchGenres(ctx, &movies[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch genres for movie %d: %w", op, movies[i].ID, err)
		}

		err = m.FetchKeywords(ctx, &movies[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch keywords for movie %d: %w", op, movies[i].ID, err)
		}

		err = m.FetchAgeCategories(ctx, &movies[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch age categories for movie %d: %w", op, movies[i].ID, err)
		}

		err = m.FetchCover(ctx, &movies[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch cover for movie %d: %w", op, movies[i].ID, err)
		}

		err = m.FetchScreenshots(ctx, &movies[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch screenshots for movie %d: %w", op, movies[i].ID, err)
		}
	}

	return movies, nil
}

func (m *MovieStorage) FetchGenres(ctx context.Context, movie *models.Movie) error {
	const op = "storage.movie.FetchGenres"

	query := `SELECT g.id, g.name 
				FROM genres g 
				JOIN movie_genres mg ON g.id = mg.genre_id 
				WHERE mg.movie_id = $1`

	rows, err := m.storage.db.QueryContext(ctx, query, movie.ID)
	if err != nil {
		return fmt.Errorf("%s: query movie genres: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var genre models.Genre // Assuming you have a Genre struct
		err := rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return fmt.Errorf("%s: scan genre row: %w", op, err)
		}
		movie.Genres = append(movie.Genres, genre)
	}

	return nil
}

func (m *MovieStorage) FetchKeywords(ctx context.Context, movie *models.Movie) error {
	const op = "storage.movie.FetchKeywords"

	query := `SELECT k.id, k.name 
				FROM key_words k 
				JOIN movie_key_words mk ON k.id = mk.key_word_id 
				WHERE mk.movie_id = $1
				`

	rows, err := m.storage.db.QueryContext(ctx, query, movie.ID)
	if err != nil {
		return fmt.Errorf("%s: query movie keywords: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var keyword models.Keyword // Assuming you have a Keyword struct
		err := rows.Scan(&keyword.ID, &keyword.Name)
		if err != nil {
			return fmt.Errorf("%s: scan keyword row: %w", op, err)
		}
		movie.Keywords = append(movie.Keywords, keyword)
	}

	return nil
}

func (m *MovieStorage) FetchMovieData(ctx context.Context) ([]models.Movie, error) {
	const op = "storage.movie.FetchMovieData"

	movies := make([]models.Movie, 0)

	rows, err := m.storage.db.QueryContext(ctx, `SELECT * FROM movies`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var movie models.Movie
		err = rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Description, &movie.Popularity, &movie.YoutubeID, &movie.Duration, &movie.Director, &movie.Producer)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		movies = append(movies, movie)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies, nil
}

func (m *MovieStorage) FetchCover(ctx context.Context, movie *models.Movie) error {
	const op = "storage.movie.FetchCover"
	var filename string
	err := m.storage.db.QueryRowContext(ctx, `SELECT filename FROM movie_covers WHERE movie_id = $1`, movie.ID).Scan(&filename)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	movie.Cover.Filename = filename
	return nil
}

func (m *MovieStorage) FetchScreenshots(ctx context.Context, movie *models.Movie) error {
	const op = "storage.movie.FetchScreenshots"

	query := "SELECT filename FROM movie_screenshots WHERE movie_id = $1"
	rows, err := m.storage.db.QueryContext(ctx, query, movie.ID)
	if err != nil {
		return fmt.Errorf("%s: query movie screenshots: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var screenshot models.Screenshot // Assuming you have a 'Screenshot' struct
		err := rows.Scan(&screenshot.Filename)
		if err != nil {
			return fmt.Errorf("%s: scan screenshot row: %w", op, err)
		}
		movie.Screenshots = append(movie.Screenshots, screenshot)
	}

	return nil
}

func (m *MovieStorage) FetchAgeCategories(ctx context.Context, movie *models.Movie) error {
	const op = "storage.movie.FetchAgeCategories"

	query := `SELECT ac.id, ac.min_age, ac.max_age
              FROM age_categories ac 
              JOIN movie_age_categories mac ON ac.id = mac.age_category_id 
              WHERE mac.movie_id = $1`

	rows, err := m.storage.db.QueryContext(ctx, query, movie.ID)
	if err != nil {
		return fmt.Errorf("%s: query age categories: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var ageCategory models.AgeCategory
		err := rows.Scan(&ageCategory.ID, &ageCategory.MinAge, &ageCategory.MaxAge)
		if err != nil {
			return fmt.Errorf("%s: scan age category row: %w", op, err)
		}
		movie.AgeCategories = append(movie.AgeCategories, ageCategory)
	}

	return nil
}

func (m *MovieStorage) GetById(ctx context.Context, id int) (models.Movie, error) {
	const op = "storage.movie.GetById"

	stmt, err := m.storage.db.Prepare(`SELECT * FROM movies WHERE id = $1`)
	if err != nil {
		return models.Movie{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var movie models.Movie
	err = stmt.QueryRowContext(ctx, id).Scan(&movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Description, &movie.Popularity, &movie.YoutubeID, &movie.Duration, &movie.Director, &movie.Producer)
	if err != nil {
		return models.Movie{}, fmt.Errorf("%s: get movie: %w", op, err)
	}

	if err = m.FetchCover(ctx, &movie); err != nil {
		return models.Movie{}, fmt.Errorf("%s: Fetch cover: %w", op, err)
	}

	if err = m.FetchScreenshots(ctx, &movie); err != nil {
		return models.Movie{}, fmt.Errorf("%s: Fetch screenshots: %w", op, err)
	}

	if err = m.FetchAgeCategories(ctx, &movie); err != nil {
		return models.Movie{}, fmt.Errorf("%s: Fetch age categories: %w", op, err)
	}

	if err = m.FetchKeywords(ctx, &movie); err != nil {
		return models.Movie{}, fmt.Errorf("%s: Fetch keywords: %w", op, err)
	}

	if err = m.FetchGenres(ctx, &movie); err != nil {
		return models.Movie{}, fmt.Errorf("%s: Fetch genres: %w", op, err)
	}

	return movie, nil
}

func (m *MovieStorage) InsertToFavorites(ctx context.Context, movieID, userID int) error {
	const op = "storage.movie.InsertToFavorites"

	tx, err := m.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	// Insert into favorites (handling potential conflicts)
	stmt, err := tx.Prepare(`INSERT INTO favorite_movies (user_id, movie_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare insert statement: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, userID, movieID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: insert into favorites: %w", op, err)
	}

	// Check if anything was actually inserted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: check rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		tx.Rollback() // No need to increment popularity if nothing was inserted
		return nil    // Or you might want to return a specific "already exists" error
	}

	// Increment popularity (only if a row was inserted)
	stmt, err = tx.Prepare(`UPDATE movies SET popularity = popularity + 1 WHERE id = $1`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare update statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, movieID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: increment popularity: %w", op, err)
	}

	// Commit if all operations succeeded
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (m *MovieStorage) DeleteFromFavorites(ctx context.Context, movieID, userID int) error {
	const op = "storage.movie.DeleteFromFavorites"

	tx, err := m.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	// Attempt to delete from favorites
	stmt, err := tx.Prepare(`DELETE FROM favorite_movies WHERE user_id = $1 AND movie_id = $2`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare delete statement: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, userID, movieID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: delete from favorites: %w", op, err)
	}

	// Check if anything was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: check rows affected: %w", op, err)
	}

	if rowsAffected == 0 {
		tx.Rollback() // No need to decrement popularity if nothing was deleted
		return nil    // Or you might want to return a specific "not found" error
	}

	// Decrement popularity (only if a row was deleted)
	stmt, err = tx.Prepare(`UPDATE movies SET popularity = popularity - 1 WHERE id = $1`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare update statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, movieID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: decrement popularity: %w", op, err)
	}

	// Commit if all operations succeeded
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (m *MovieStorage) GetFavorites(ctx context.Context, userID int) ([]models.Movie, error) {
	const op = "storage.movie.GetFavorites"

	var movieIDs []int
	stmt, err := m.storage.db.Prepare(`SELECT movie_id FROM favorite_movies WHERE user_id = $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement (phase 1): %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: query (phase 1): %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var movieID int
		err = rows.Scan(&movieID)
		if err != nil {
			return nil, fmt.Errorf("%s: scan movie ID (phase 1): %w", op, err)
		}
		movieIDs = append(movieIDs, movieID)
	}

	if len(movieIDs) == 0 {
		return []models.Movie{}, nil
	}

	placeholders := make([]string, len(movieIDs))
	args := make([]interface{}, len(movieIDs))
	for i, id := range movieIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT * FROM movies WHERE id IN (%s)`, strings.Join(placeholders, ", "))

	rows, err = m.storage.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query (phase 2): %w", op, err)
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err = rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Description,
			&movie.Popularity,
			&movie.YoutubeID,
			&movie.Duration,
			&movie.Director,
			&movie.Producer,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan movie (phase 2): %w", op, err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieStorage) GetByTitle(ctx context.Context, title string) ([]models.Movie, error) {
	const op = "storage.movie.GetByTitle"

	stmt, err := m.storage.db.Prepare(`SELECT * FROM movies WHERE title ILIKE $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, "%"+title+"%")
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var movie models.Movie
		err = rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Description,
			&movie.Popularity,
			&movie.YoutubeID,
			&movie.Duration,
			&movie.Director,
			&movie.Producer,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan movie: %w", op, err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}

func (m *MovieStorage) GetByGenres(ctx context.Context, genres []string) ([]models.Movie, error) {
	const op = "storage.movie.GetByGenres"

	// Start building the SQL query
	query := `
        SELECT m.id, m.title, m.release_year, m.description, m.popularity, m.youtube_id, 
               m.duration, m.director, m.producer
        FROM movies m
        JOIN movie_genres mg ON m.id = mg.movie_id
        JOIN genres g ON mg.genre_id = g.id
        WHERE g.name IN (`

	// Add placeholders for genre names
	placeholders := make([]string, len(genres))
	for i := range genres {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	query += strings.Join(placeholders, ", ") + ")"

	// Convert []string to []interface{}
	args := make([]interface{}, len(genres))
	for i, v := range genres {
		args[i] = interface{}(v)
	}

	// Execute the query
	rows, err := m.storage.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query movies by genres: %w", op, err)
	}
	defer rows.Close()

	// Process the results
	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		if err := rows.Scan(
			&movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Description,
			&movie.Popularity, &movie.YoutubeID, &movie.Duration, &movie.Director, &movie.Producer,
		); err != nil {
			return nil, fmt.Errorf("%s: scan movie row: %w", op, err)
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iterate movie rows: %w", op, err)
	}

	return movies, nil
}

func (m *MovieStorage) GetByYear(ctx context.Context, yearStart, yearEnd int) ([]models.Movie, error) {
	const op = "storage.movie.GetByYear"

	stmt, err := m.storage.db.Prepare(`SELECT * FROM movies WHERE release_year BETWEEN $1 AND $2`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, yearStart, yearEnd)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var movie models.Movie
		err = rows.Scan(&movie.ID, &movie.Title, &movie.ReleaseYear, &movie.Description, &movie.Popularity, &movie.YoutubeID, &movie.Duration, &movie.Director, &movie.Producer)
		if err != nil {
			return nil, fmt.Errorf("%s: scan movie: %w", op, err)
		}
		movies = append(movies, movie)
	}

	return movies, nil
}
