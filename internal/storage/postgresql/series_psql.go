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

type SeriesStorage struct {
	storage *Postgres
}

func NewSeriesStorage(db *Postgres) *SeriesStorage {
	return &SeriesStorage{storage: db}
}

func (s *SeriesStorage) Insert(ctx context.Context, series models.Series) (int, error) {
	const op = "storage.series.Insert"

	tx, err := s.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	var seriesID int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO series (title, release_year, description, popularity, duration, director, producer) 
         VALUES ($1, $2, $3, $4, $5, $6, $7) 
         RETURNING id`,
		series.Title, series.ReleaseYear, series.Description, series.Popularity,
		series.Duration, series.Director, series.Producer,
	).Scan(&seriesID)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: insert series: %w", op, err)
	}

	// Insert Genres
	for _, genre := range series.Genres {
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

		// 3. Link Genre with Series
		_, err = tx.ExecContext(ctx,
			`INSERT INTO series_genres (series_id, genre_id) 
			 VALUES ($1, $2)`,
			seriesID, genreID,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: link genre to series: %w", op, err)
		}
	}

	// Insert Keywords
	for _, keyword := range series.Keywords {
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

		// 3. Link Keyword with Series
		_, err = tx.ExecContext(ctx,
			`INSERT INTO series_key_words (series_id, key_word_id) 
		 VALUES ($1, $2)`,
			seriesID, keywordID,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: link keyword to series: %w", op, err)
		}
	}

	// Insert Age Categories
	for _, ageCategory := range series.AgeCategories {
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

		// 3. Link Age Category with Series
		_, err = tx.ExecContext(ctx,
			`INSERT INTO series_age_categories (series_id, age_category_id) 
		 VALUES ($1, $2)`,
			seriesID, ageCategoryID,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: link age category to series: %w", op, err)
		}
	}
	// Insert Screenshots
	for _, screenshot := range series.Screenshots {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO series_screenshots (series_id, filename) 
			VALUES ($1, $2)`,
			seriesID, screenshot.Filename,
		)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: insert screenshots: %w", op, err)
		}
	}

	// Insert Cover
	_, err = tx.ExecContext(ctx,
		`INSERT INTO series_covers (series_id, filename) 
         VALUES ($1, $2)`,
		seriesID, series.Cover.Filename,
	)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: insert cover: %w", op, err)
	}

	// 7. Insert Seasons and Episodes
	for _, season := range series.Seasons {
		// Insert season
		var seasonID int
		err = tx.QueryRowContext(ctx,
			`INSERT INTO seasons (series_id, season_number) 
                 VALUES ($1, $2)
                 RETURNING id`,
			seriesID, season.SeasonNumber).Scan(&seasonID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: insert season: %w", op, err)
		}

		// Insert episodes for this season
		for _, episode := range season.Episodes {
			_, err = tx.ExecContext(ctx,
				`INSERT INTO episodes (season_id, episode_number, youtube_id)
				VALUES ($1, $2, $3)`,
				seasonID, episode.EpisodeNumber, episode.Link)
			if err != nil {
				var pqErr *pq.Error
				if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
					return 0, fmt.Errorf("%s: %w", op, storage.ErrSeriesExists)
				}
				tx.Rollback()
				return 0, fmt.Errorf("%s: insert episode: %w", op, err)
			}
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return seriesID, nil
}

func (s *SeriesStorage) Update(ctx context.Context, series models.Series) error {
	const op = "storage.series.Update"

	tx, err := s.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	// Update Series (Similar to Series Update)
	if series.Title != "" || series.ReleaseYear != 0 || series.Description != "" ||
		series.Popularity != 0 || series.Duration != 0 ||
		series.Director != "" || series.Producer != "" {
		_, err := tx.ExecContext(ctx,
			`UPDATE series 
             SET title = COALESCE($1, title), release_year = COALESCE($2, release_year), 
             description = COALESCE($3, description), popularity = COALESCE($4, popularity), 
             duration = COALESCE($5, duration), director = COALESCE($6, director), 
             producer = COALESCE($7, producer)
             WHERE id = $8`,
			series.Title, series.ReleaseYear, series.Description, series.Popularity,
			series.Duration, series.Director, series.Producer,
			series.ID,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: update series: %w", op, err)
		}
	}

	// Update Age Categories
	if len(series.AgeCategories) > 0 {
		// Delete existing
		_, err := tx.ExecContext(ctx, "DELETE FROM series_age_categories WHERE series_id = $1", series.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: delete existing age categories: %w", op, err)
		}

		// Insert new age categories
		for _, ageCategory := range series.AgeCategories {
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

			var ageCategoryID int
			err = tx.QueryRowContext(ctx,
				`SELECT id FROM age_categories WHERE min_age = $1 AND max_age = $2`,
				ageCategory.MinAge, ageCategory.MaxAge,
			).Scan(&ageCategoryID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: get age category id: %w", op, err)
			}

			_, err = tx.ExecContext(ctx,
				`INSERT INTO series_age_categories (series_id, age_category_id) 
                 VALUES ($1, $2)`,
				series.ID, ageCategoryID,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: link age category to series: %w", op, err)
			}
		}
	}

	// Update Genres
	if len(series.Genres) > 0 {
		// Delete existing genres
		_, err := tx.ExecContext(ctx, "DELETE FROM series_genres WHERE series_id = $1", series.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: delete existing genres: %w", op, err)
		}

		// Insert new genres
		for _, genre := range series.Genres {
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

			var genreID int
			err = tx.QueryRowContext(ctx,
				`SELECT id FROM genres WHERE name = $1`,
				genre.Name,
			).Scan(&genreID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: get genre id: %w", op, err)
			}

			_, err = tx.ExecContext(ctx,
				`INSERT INTO series_genres (series_id, genre_id) 
             VALUES ($1, $2)`,
				series.ID, genreID,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: link genre to series: %w", op, err)
			}
		}
	}

	// Update Keywords
	if len(series.Keywords) > 0 {
		// Delete existing keywords
		_, err := tx.ExecContext(ctx, "DELETE FROM series_key_words WHERE series_id = $1", series.ID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: delete existing keywords: %w", op, err)
		}

		// Insert new keywords
		for _, keyword := range series.Keywords {
			_, err := tx.ExecContext(ctx,
				`INSERT INTO key_words (name) 
             VALUES ($1)
             ON CONFLICT (name) DO NOTHING`,
				keyword.Name,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: insert keyword: %w", op, err)
			}

			var keywordID int
			err = tx.QueryRowContext(ctx,
				`SELECT id FROM key_words WHERE name = $1`,
				keyword.Name,
			).Scan(&keywordID)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: get keyword id: %w", op, err)
			}

			_, err = tx.ExecContext(ctx,
				`INSERT INTO series_key_words (series_id, key_word_id) 
             VALUES ($1, $2)`,
				series.ID, keywordID,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: link keyword to series: %w", op, err)
			}
		}
	}

	// Update Seasons and Episodes
	for _, season := range series.Seasons {
		row := tx.QueryRowContext(ctx,
			`SELECT id, season_number FROM seasons WHERE series_id = $1 AND season_number = $2`,
			series.ID, season.SeasonNumber,
		)

		err := row.Scan(&season.ID, &season.SeasonNumber)
		if err != nil {
			fmt.Println("err: ", err)
		}
		// if err != nil {
		// 	if err == sql.ErrNoRows {
		// 		// Insert new Season
		// 		_, err = tx.ExecContext(ctx,
		// 			`INSERT INTO seasons (series_id, season_number)
		// 			VALUES ($1, $2)`,
		// 			series.ID, season.SeasonNumber,
		// 		)
		// 		if err != nil {
		// 			tx.Rollback()
		// 			return fmt.Errorf("%s: insert season: %w", op, err)
		// 		}
		// 	} else {
		// 		tx.Rollback()
		// 		return fmt.Errorf("%s: get season id: %w", op, err)
		// 	}
		// }

		_, err = tx.ExecContext(ctx,
			`UPDATE seasons
			SET season_number = $1
			WHERE id = $2`,
			season.SeasonNumber, season.ID,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: update season: %w", op, err)
		}

		// Update Episodes for this Season
		for _, episode := range season.Episodes {

			row := tx.QueryRowContext(ctx,
				`SELECT id, episode_number FROM episodes WHERE season_id = $1 AND episode_number = $2`,
				season.ID, episode.EpisodeNumber,
			)

			err := row.Scan(&episode.ID, &episode.EpisodeNumber)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: get episode id: %w", op, err)
			}
			_, err = tx.ExecContext(ctx,
				`UPDATE episodes
				SET episode_number = $1, youtube_id = $2
				WHERE id = $3`,
				episode.EpisodeNumber, episode.Link, episode.ID,
			)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("%s: update episode: %w", op, err)
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

func (s *SeriesStorage) UpdateCover(ctx context.Context, seriesID int, cover models.Cover) error {
	const op = "storage.series.UpdateCover"

	tx, err := s.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// 1. Delete existing cover (if any)
	_, err = tx.ExecContext(ctx, `DELETE FROM series_covers WHERE series_id = $1`, seriesID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: delete error: %w", op, err)
	}

	// 2. Insert the new cover
	_, err = tx.ExecContext(ctx, `
        INSERT INTO series_covers (series_id, filename) VALUES ($1, $2)`, seriesID, cover.Filename)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: insert error: %w", op, err)
	}

	return tx.Commit()
}

func (s *SeriesStorage) UpdateScreenshots(ctx context.Context, seriesID int, screenshots []models.Screenshot) error {
	const op = "storage.series.UpdateScreenshots"

	tx, err := s.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// 1. Delete existing screenshots
	_, err = tx.ExecContext(ctx, `DELETE FROM series_screenshots WHERE series_id = $1`, seriesID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: delete error: %w", op, err)
	}

	// 2. Insert new screenshots
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO series_screenshots (series_id, filename) VALUES ($1, $2)`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare error: %w", op, err)
	}
	defer stmt.Close()

	for _, screenshot := range screenshots {
		_, err = stmt.ExecContext(ctx, seriesID, screenshot.Filename)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("%s: insert error: %w", op, err)
		}
	}

	return tx.Commit()
}

func (s *SeriesStorage) Delete(ctx context.Context, id int) error {
	const op = "storage.series.Delete"

	stmt, err := s.storage.db.Prepare(`DELETE FROM series WHERE id = $1`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SeriesStorage) GetById(ctx context.Context, id int) (models.Series, error) {
	const op = "storage.series.GetById"

	stmt, err := s.storage.db.Prepare(`SELECT * FROM series WHERE id = $1`)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var series models.Series

	err = stmt.QueryRowContext(ctx, id).Scan(&series.ID, &series.Title, &series.ReleaseYear, &series.Description, &series.Popularity, &series.Duration, &series.Director, &series.Producer)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: %w", op, err)
	}

	err = s.FetchGenres(ctx, &series)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: Fetch genres for series %d: %w", op, series.ID, err)
	}

	err = s.FetchKeywords(ctx, &series)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: Fetch keywords for series %d: %w", op, series.ID, err)
	}

	err = s.FetchCover(ctx, &series)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: Fetch covers for series %d: %w", op, series.ID, err)
	}

	err = s.FetchScreenshots(ctx, &series)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: Fetch screenshots for series %d: %w", op, series.ID, err)
	}

	err = s.FetchSeasons(ctx, &series)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: Fetch seasons for series %d: %w", op, series.ID, err)
	}

	err = s.FetchAgeCategories(ctx, &series)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: Fetch age categories for series %d: %w", op, series.ID, err)
	}

	return series, nil
}

func (s *SeriesStorage) GetAll(ctx context.Context) ([]models.Series, error) {
	const op = "storage.series.GetAll"

	series, err := s.FetchSeriesData(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for i := range series {
		err = s.FetchGenres(ctx, &series[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch genres for series %d: %w", op, series[i].ID, err)
		}

		err = s.FetchKeywords(ctx, &series[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch keywords for series %d: %w", op, series[i].ID, err)
		}

		err = s.FetchCover(ctx, &series[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch covers for series %d: %w", op, series[i].ID, err)
		}

		err = s.FetchSeasons(ctx, &series[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch seasons for series %d: %w", op, series[i].ID, err)
		}

		err = s.FetchAgeCategories(ctx, &series[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch age categories for series %d: %w", op, series[i].ID, err)
		}

		err = s.FetchScreenshots(ctx, &series[i])
		if err != nil {
			return nil, fmt.Errorf("%s: Fetch screenshots for series %d: %w", op, series[i].ID, err)
		}
	}

	return series, nil
}

// func (s *SeriesStorage) GetSeasons(ctx context.Context, seriesID int) ([]models.Season, error) {
// 	const op = "storage.series.GetSeasons"

// 	seasons, err := s.FetchSeasonsData(ctx, seriesID)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: %w", op, err)
// 	}

// 	for i := range seasons {
// 		err = s.FetchEpisodesData(ctx, seriesID, seasons[i].ID)
// 		if err != nil {
// 			return nil, fmt.Errorf("%s: Fetch episodes for season %d: %w", op, seasons[i].ID, err)
// 		}
// 	}

// 	return seasons, nil
// }

func (s *SeriesStorage) FetchSeasons(ctx context.Context, series *models.Series) error {
	const op = "storage.series.FetchSeasons"

	query := `SELECT s.id, s.season_number
				FROM seasons s
				WHERE s.series_id = $1`

	rows, err := s.storage.db.QueryContext(ctx, query, series.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	seasons := make([]models.Season, 0)

	for rows.Next() {
		var season models.Season
		err = rows.Scan(&season.ID, &season.SeasonNumber)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		season.Episodes, err = s.FetchEpisodes(ctx, series.ID, season.ID)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		season.SeriesID = series.ID

		seasons = append(seasons, season)
	}

	series.Seasons = seasons

	return nil
}

func (s *SeriesStorage) FetchEpisodes(ctx context.Context, seriesID, seasonID int) ([]models.Episode, error) {
	const op = "storage.series.FetchEpisodes"

	query := `SELECT id, episode_number, youtube_id
              FROM episodes 
              WHERE season_id = $1 AND season_id IN (SELECT id FROM seasons WHERE series_id = $2)`

	rows, err := s.storage.db.QueryContext(ctx, query, seasonID, seriesID) // Corrected order
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var episodes []models.Episode
	for rows.Next() {
		var episode models.Episode
		err := rows.Scan(&episode.ID, &episode.EpisodeNumber, &episode.Link)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		episode.SeasonID = seasonID
		episodes = append(episodes, episode)
	}

	return episodes, nil
}

func (s *SeriesStorage) FetchGenres(ctx context.Context, series *models.Series) error {
	const op = "storage.series.FetchGenres"

	query := `SELECT g.id, g.name 
				FROM genres g 
				JOIN series_genres sg ON sg.genre_id = g.id
				WHERE sg.series_id = $1`

	rows, err := s.storage.db.QueryContext(ctx, query, series.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	genres := make([]models.Genre, 0)

	for rows.Next() {
		var genre models.Genre
		err = rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		genres = append(genres, genre)
	}

	series.Genres = genres

	return nil
}

func (s *SeriesStorage) FetchKeywords(ctx context.Context, series *models.Series) error {
	const op = "storage.series.FetchKeywords"

	query := `SELECT k.id, k.name
				FROM key_words k
				JOIN series_key_words sk ON sk.key_word_id = k.id
				WHERE sk.series_id = $1`

	rows, err := s.storage.db.QueryContext(ctx, query, series.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	keywords := make([]models.Keyword, 0)

	for rows.Next() {
		var keyword models.Keyword
		err = rows.Scan(&keyword.ID, &keyword.Name)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		keywords = append(keywords, keyword)
	}

	series.Keywords = keywords

	return nil
}

func (s *SeriesStorage) FetchSeriesData(ctx context.Context) ([]models.Series, error) {
	const op = "storage.series.FetchSeriesData"

	series := make([]models.Series, 0)

	rows, err := s.storage.db.QueryContext(ctx, `SELECT * FROM series`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var series_data models.Series
		err = rows.Scan(&series_data.ID, &series_data.Title, &series_data.ReleaseYear, &series_data.Description, &series_data.Popularity, &series_data.Duration, &series_data.Director, &series_data.Producer)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		series = append(series, series_data)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return series, nil
}

func (s *SeriesStorage) FetchCover(ctx context.Context, series *models.Series) error {
	const op = "storage.series.FetchCover"
	var filename string
	err := s.storage.db.QueryRowContext(ctx, `SELECT filename FROM series_covers WHERE series_id = $1`, series.ID).Scan(&filename)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	series.Cover.Filename = filename
	return nil
}

func (s *SeriesStorage) FetchScreenshots(ctx context.Context, series *models.Series) error {
	const op = "storage.series.FetchScreenshots"

	query := "SELECT filename FROM series_screenshots WHERE series_id = $1"
	rows, err := s.storage.db.QueryContext(ctx, query, series.ID)
	if err != nil {
		return fmt.Errorf("%s: query series screenshots: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var screenshot models.Screenshot // Assuming you have a 'Screenshot' struct
		err := rows.Scan(&screenshot.Filename)
		if err != nil {
			return fmt.Errorf("%s: scan screenshot row: %w", op, err)
		}
		series.Screenshots = append(series.Screenshots, screenshot)
	}

	return nil
}

func (s *SeriesStorage) FetchAgeCategories(ctx context.Context, series *models.Series) error {
	const op = "storage.series.FetchAgeCategories"

	query := `SELECT ac.id, ac.min_age, ac.max_age
              FROM age_categories ac 
              JOIN series_age_categories sac ON ac.id = sac.age_category_id 
              WHERE sac.series_id = $1`

	rows, err := s.storage.db.QueryContext(ctx, query, series.ID)
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
		series.AgeCategories = append(series.AgeCategories, ageCategory)
	}

	return nil
}

func (s *SeriesStorage) InsertToFavorites(ctx context.Context, seriesID, userID int) error {
	const op = "storage.series.InsertToFavorites"

	tx, err := s.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	// Insert into favorites (handling potential conflicts)
	stmt, err := tx.Prepare(`INSERT INTO favorite_series (series_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare insert statement: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, seriesID, userID)
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
		tx.Rollback()
		return nil
	}

	// Increment popularity (only if a row was inserted)
	stmt, err = tx.Prepare(`UPDATE series SET popularity = popularity + 1 WHERE id = $1`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare update statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, seriesID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: increment popularity: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (s *SeriesStorage) DeleteFromFavorites(ctx context.Context, seriesID, userID int) error {
	const op = "storage.series.DeleteFromFavorites"

	tx, err := s.storage.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	// Attempt to delete from favorites
	stmt, err := tx.Prepare(`DELETE FROM favorite_series WHERE series_id = $1 AND user_id = $2`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare delete statement: %w", op, err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, seriesID, userID)
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
		tx.Rollback()
		return nil // Or a specific "not found" error if desired
	}

	// Decrement popularity (only if a row was deleted)
	stmt, err = tx.Prepare(`UPDATE series SET popularity = popularity - 1 WHERE id = $1`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: prepare update statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, seriesID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("%s: decrement popularity: %w", op, err)
	}

	// Commit if a favorite was deleted
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}

func (s *SeriesStorage) GetFavorites(ctx context.Context, userID int) ([]models.Series, error) {
	const op = "storage.series.GetFavorites"

	// Phase 1: Get favorite series IDs
	var seriesIDs []int
	stmt, err := s.storage.db.Prepare(`SELECT series_id FROM favorite_series WHERE user_id = $1`)
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
		var seriesID int
		err = rows.Scan(&seriesID)
		if err != nil {
			return nil, fmt.Errorf("%s: scan series ID (phase 1): %w", op, err)
		}
		seriesIDs = append(seriesIDs, seriesID)
	}

	// Phase 2: Get series details if favorites exist
	if len(seriesIDs) == 0 {
		return []models.Series{}, nil
	}

	placeholders := make([]string, len(seriesIDs))
	args := make([]interface{}, len(seriesIDs))
	for i, id := range seriesIDs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}

	query := fmt.Sprintf(`SELECT * FROM series WHERE id IN (%s)`, strings.Join(placeholders, ", "))

	rows, err = s.storage.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query (phase 2): %w", op, err)
	}
	defer rows.Close()

	var series []models.Series
	for rows.Next() {
		var s models.Series
		err = rows.Scan(
			&s.ID,
			&s.Title,
			&s.ReleaseYear,
			&s.Description,
			&s.Popularity,
			&s.Duration,
			&s.Director,
			&s.Producer,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: scan series (phase 2): %w", op, err)
		}

		series = append(series, s)
	}

	return series, nil
}

func (s *SeriesStorage) GetByTitle(ctx context.Context, title string) ([]models.Series, error) {
	const op = "storage.series.GetByTitle"

	stmt, err := s.storage.db.Prepare(`SELECT * FROM series WHERE title ILIKE $1`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, "%"+title+"%")
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var series []models.Series
	for rows.Next() {
		var s models.Series
		err = rows.Scan(&s.ID, &s.Title, &s.ReleaseYear, &s.Description, &s.Popularity, &s.Duration, &s.Director, &s.Producer)
		if err != nil {
			return nil, fmt.Errorf("%s: scan series: %w", op, err)
		}
		series = append(series, s)
	}

	return series, nil
}

func (s *SeriesStorage) GetByGenres(ctx context.Context, genres []string) ([]models.Series, error) {
	const op = "storage.series.GetByGenres"

	// Start building the SQL query
	query := `
        SELECT s.id, s.title, s.release_year, s.description, s.popularity, 
               s.duration, s.director, s.producer
        FROM series s
        JOIN series_genres sg ON s.id = sg.series_id
        JOIN genres g ON sg.genre_id = g.id
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
	rows, err := s.storage.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: query series by genres: %w", op, err)
	}
	defer rows.Close()

	// Process the results
	var seriesList []models.Series
	for rows.Next() {
		var series models.Series
		if err := rows.Scan(
			&series.ID, &series.Title, &series.ReleaseYear, &series.Description,
			&series.Popularity, &series.Duration, &series.Director, &series.Producer,
		); err != nil {
			return nil, fmt.Errorf("%s: scan series row: %w", op, err)
		}
		seriesList = append(seriesList, series)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: iterate series rows: %w", op, err)
	}

	return seriesList, nil
}

func (s *SeriesStorage) GetByYear(ctx context.Context, yearStart, yearEnd int) ([]models.Series, error) {
	const op = "storage.series.GetByYear"

	stmt, err := s.storage.db.Prepare(`SELECT * FROM series WHERE release_year BETWEEN $1 AND $2`)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, yearStart, yearEnd)
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var series []models.Series
	for rows.Next() {
		var s models.Series
		err = rows.Scan(&s.ID, &s.Title, &s.ReleaseYear, &s.Description, &s.Popularity, &s.Duration, &s.Director, &s.Producer)
		if err != nil {
			return nil, fmt.Errorf("%s: scan series: %w", op, err)
		}
		series = append(series, s)
	}

	return series, nil
}
