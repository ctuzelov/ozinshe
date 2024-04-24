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

type SeriesService struct {
	Storage psql.Series
}

func NewSeriesService(series psql.Series) *SeriesService {
	return &SeriesService{Storage: series}
}

func (s *SeriesService) Add(series models.Series, image_data models.SavePhoto) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	const op = "service.series.Add"

	var series_screenshots []models.Screenshot

	screenshots := image_data.File_form.File["screenshots"]
	for _, fileHeader := range screenshots {
		if err := validation.ValidateImageFile(fileHeader, image_data.MaxImageSize); err != nil {
			return 0, err
		}

		series_screenshots = append(series_screenshots, models.Screenshot{
			Filename: fileHeader.Filename,
		})
	}

	cover := image_data.File_form.File["cover"]
	if err := validation.ValidateImageFile(cover[0], image_data.MaxImageSize); err != nil {
		return 0, err
	}

	series_cover := models.Cover{
		Filename: cover[0].Filename,
	}

	series.Screenshots = series_screenshots
	series.Cover = series_cover

	series.ID, err = s.Storage.Insert(ctx, series)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	uploadPath := image_data.UploadPath + strconv.Itoa(series.ID) + "/covers/" + cover[0].Filename
	err = helper.ProcessSaving(cover[0], uploadPath)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	for _, fileHeader := range screenshots {
		uploadPath := image_data.UploadPath + strconv.Itoa(series.ID) + "/screenshots/" + fileHeader.Filename

		err = helper.ProcessSaving(fileHeader, uploadPath)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
		}
	}

	return series.ID, nil
}

func (s *SeriesService) GetAll() ([]models.Series, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.GetAll"

	series, err := s.Storage.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return series, nil
}

func (s *SeriesService) GetById(id int) (models.Series, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.GetById"

	series, err := s.Storage.GetById(ctx, id)
	if err != nil {
		return models.Series{}, fmt.Errorf("%s: %w", op, err)
	}

	return series, nil
}

func (s *SeriesService) GetSeason(seriesID, seasonNumber int) ([]models.Episode, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.GetSeason"

	episodes, err := s.Storage.FetchEpisodes(ctx, seriesID, seasonNumber)
	if err != nil {
		return []models.Episode{}, fmt.Errorf("%s: %w", op, err)
	}

	return episodes, nil
}

func (s *SeriesService) Remove(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.Remove"

	err := s.Storage.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = helper.DeleteDirectory("uploads/series/" + strconv.Itoa(id))
	if err != nil {
		return fmt.Errorf("%s: delete directory: %w", op, err)
	}

	return nil
}

func (s *SeriesService) Update(id int, series models.Series) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.Update"

	err := s.Storage.Update(ctx, series)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SeriesService) UpdateCover(id int, image_data models.SavePhoto) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.UpdateCover"

	err := helper.DeleteDirectory("uploads/series/" + strconv.Itoa(id) + "/covers")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	cover := image_data.File_form.File["cover"]
	if err := validation.ValidateImageFile(cover[0], image_data.MaxImageSize); err != nil {
		return err
	}

	err = s.Storage.UpdateCover(ctx, id, models.Cover{
		Filename: cover[0].Filename,
	})

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	uploadPath := image_data.UploadPath + strconv.Itoa(id) + "/covers/" + cover[0].Filename
	err = helper.ProcessSaving(cover[0], uploadPath)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SeriesService) UpdateScreenshots(id int, image_data models.SavePhoto) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.UpdateScreenshots"

	err := helper.DeleteDirectory("uploads/series/" + strconv.Itoa(id) + "/screenshots")

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	var series_screenshots []models.Screenshot

	screenshots := image_data.File_form.File["screenshots"]
	for _, fileHeader := range screenshots {
		if err := validation.ValidateImageFile(fileHeader, image_data.MaxImageSize); err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}

		series_screenshots = append(series_screenshots, models.Screenshot{
			Filename: fileHeader.Filename,
		})
	}

	err = s.Storage.UpdateScreenshots(ctx, id, series_screenshots)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, fileHeader := range screenshots {
		uploadPath := image_data.UploadPath + strconv.Itoa(id) + "/screenshots/" + fileHeader.Filename

		err = helper.ProcessSaving(fileHeader, uploadPath)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (s *SeriesService) AddToFavorites(seriesID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.AddToFavorites"

	err := s.Storage.InsertToFavorites(ctx, seriesID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SeriesService) RemoveFromFavorites(seriesID, userID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.RemoveFromFavorites"

	err := s.Storage.DeleteFromFavorites(ctx, seriesID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *SeriesService) GetFavorites(userID int) ([]models.Series, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.GetFavorites"

	series_list, err := s.Storage.GetFavorites(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return series_list, nil
}

func (series *SeriesService) GetFiltered(filter models.FilterParams) ([]models.Series, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	const op = "service.series.GetFiltered"

	var err error
	var seriesWithTitle []models.Series
	if filter.Title != "" {
		seriesWithTitle, err = series.Storage.GetByTitle(ctx, filter.Title)
		if err != nil {
			return nil, fmt.Errorf("%s: error fetching series by title: %w", op, err)
		}
	}

	if len(seriesWithTitle) == 1 {
		return seriesWithTitle, nil
	}

	var seriesWithGenres []models.Series
	if len(filter.Genres) != 0 {
		seriesWithGenres, err = series.Storage.GetByGenres(ctx, filter.Genres)
		if err != nil {
			return nil, fmt.Errorf("%s: error fetching series by genres: %w", op, err)
		}
	}

	var seriesWithYear []models.Series
	if filter.YearStart != 0 || filter.YearEnd != 0 {
		for _, series := range seriesWithGenres {
			if series.ReleaseYear >= filter.YearStart && series.ReleaseYear <= filter.YearEnd {
				seriesWithYear = append(seriesWithYear, series)
			}
		}
	}

	if err != nil {
		return nil, fmt.Errorf("%s: error during final filtering: %w", op, err)
	}

	return seriesWithYear, nil
}
