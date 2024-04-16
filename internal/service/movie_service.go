package service

import (
	"context"
	"fmt"
	"ozinshe/internal/helper"
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
	"ozinshe/internal/validation"
	"time"
)

type MovieService struct {
	Storage     psql.Movie
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

		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
		uploadPath := image_data.UploadPath + "screenshots/" + fileHeader.Filename

		err = helper.ProcessSaving(fileHeader, uploadPath)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}

		movie_screenshots = append(movie_screenshots, models.Screenshot{
			Filename: fileHeader.Filename,
		})
	}

	cover := image_data.File_form.File["cover"]
	if err := validation.ValidateImageFile(cover[0], image_data.MaxImageSize); err != nil {
		return 0, err
	}

	uploadPath := image_data.UploadPath + "cover/" + cover[0].Filename
	err = helper.ProcessSaving(cover[0], uploadPath)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
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

	return movie.ID, nil
}
