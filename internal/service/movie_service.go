package service

import (
	"context"
	"fmt"
	"ozinshe/internal/helper"
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
	"ozinshe/internal/validation"
	"strconv"
	"time"
)

type MovieService struct {
	Storage psql.Movie
}

func NewMovieService(storage psql.Movie) *MovieService {
	return &MovieService{Storage: storage}
}

func (movies *MovieService) Add(movie models.Movie, image_data models.SavePhoto) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	const op = "service.movie.Add"

	var movie_screenshots []models.Screenshot

	screenshots := image_data.File_form.File["screenshots"]
	for _, fileHeader := range screenshots {
		if err := validation.ValidateImageFile(fileHeader, image_data.MaxImageSize); err != nil {
			return 0, err
		}

		movie_screenshots = append(movie_screenshots, models.Screenshot{
			Filename: fileHeader.Filename,
		})
	}

	cover := image_data.File_form.File["cover"]
	if err := validation.ValidateImageFile(cover[0], image_data.MaxImageSize); err != nil {
		return 0, err
	}

	movie_cover := models.Cover{
		Filename: cover[0].Filename,
	}

	movie.Screenshots = movie_screenshots
	movie.Cover = movie_cover

	movie.ID, err = movies.Storage.Insert(ctx, movie)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	uploadPath := image_data.UploadPath + "/" + strconv.Itoa(movie.ID) + "/covers/" + cover[0].Filename
	err = helper.ProcessSaving(cover[0], uploadPath)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	for _, fileHeader := range screenshots {
		uploadPath := image_data.UploadPath + "/" + strconv.Itoa(movie.ID) + "/screenshots/" + fileHeader.Filename

		err = helper.ProcessSaving(fileHeader, uploadPath)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	return movie.ID, nil
}

func (movies *MovieService) Remove(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.Remove"
	err := movies.Storage.Delete(ctx, id)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = helper.DeleteDirectory("uploads/movies/" + strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (movies *MovieService) Update(id int, movie models.Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.Update"

	movie.ID = id
	err := movies.Storage.Update(ctx, movie)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (movies *MovieService) UpdateCover(id int, image_data models.SavePhoto) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.UpdateCover"

	err := helper.DeleteDirectory("uploads/movies/" + strconv.Itoa(id) + "/covers")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	cover := image_data.File_form.File["cover"]
	if err := validation.ValidateImageFile(cover[0], image_data.MaxImageSize); err != nil {
		return err
	}

	uploadPath := image_data.UploadPath + "/" + strconv.Itoa(id) + "/covers/" + cover[0].Filename
	err = helper.ProcessSaving(cover[0], uploadPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = movies.Storage.UpdateCover(ctx, id, models.Cover{
		Filename: cover[0].Filename,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (movies *MovieService) UpdateScreenshots(id int, image_data models.SavePhoto) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.UpdateScreenshots"
	var err error

	err = helper.DeleteDirectory("uploads/movies/" + strconv.Itoa(id) + "/screenshots")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var movie_screenshots []models.Screenshot

	screenshots := image_data.File_form.File["screenshots"]
	for _, fileHeader := range screenshots {
		if err := validation.ValidateImageFile(fileHeader, image_data.MaxImageSize); err != nil {
			return err
		}

		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		uploadPath := image_data.UploadPath + "/" + strconv.Itoa(id) + "/screenshots/" + fileHeader.Filename

		err = helper.ProcessSaving(fileHeader, uploadPath)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		movie_screenshots = append(movie_screenshots, models.Screenshot{
			Filename: fileHeader.Filename,
		})
	}

	err = movies.Storage.UpdateScreenshots(ctx, id, movie_screenshots)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (movies *MovieService) GetById(id int) (models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.GetById"

	movie, err := movies.Storage.GetById(ctx, id)
	if err != nil {
		return models.Movie{}, fmt.Errorf("%s: %w", op, err)
	}

	return movie, nil
}

func (movies *MovieService) GetAll() ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.GetAll"

	movies_list, err := movies.Storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies_list, nil
}

func (movies *MovieService) AddToFavorites(movieID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.AddToFavorites"

	err := movies.Storage.InsertToFavorites(ctx, movieID, userID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (movies *MovieService) RemoveFromFavorites(movieID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.RemoveFromFavorites"

	err := movies.Storage.DeleteFromFavorites(ctx, movieID, userID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (movies *MovieService) GetFavorites(userID int) ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.GetFavorites"

	movies_list, err := movies.Storage.GetFavorites(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return movies_list, nil
}

func (movies *MovieService) GetFiltered(filter models.FilterParams) ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.movie.GetFiltered"

	var err error
	var moviesWithTitle []models.Movie
	if filter.Title != "" {
		moviesWithTitle, err = movies.Storage.GetByTitle(ctx, filter.Title)
		if err != nil {
			return nil, fmt.Errorf("%s: error fetching movies by title: %w", op, err)
		}
	}

	if len(moviesWithTitle) == 1 {
		return moviesWithTitle, nil
	}

	var moviesWithGenres []models.Movie
	if len(filter.Genres) != 0 {
		moviesWithGenres, err = movies.Storage.GetByGenres(ctx, filter.Genres)
		if err != nil {
			return nil, fmt.Errorf("%s: error fetching movies by genres: %w", op, err)
		}
	}

	var moviesWithYear []models.Movie
	if filter.YearStart != 0 || filter.YearEnd != 0 {
		if len(moviesWithGenres) == 0{
			moviesWithYear, err = movies.Storage.GetByYear(ctx, filter.YearStart, filter.YearEnd)
			if err != nil {
				return nil, fmt.Errorf("%s: error fetching movies by year: %w", op, err)
			}
		}else{
			for _, movie := range moviesWithGenres {
				if movie.ReleaseYear >= filter.YearStart && movie.ReleaseYear <= filter.YearEnd {
					moviesWithYear = append(moviesWithYear, movie)
				}
			}
		}
	}

	return moviesWithYear, nil
}
