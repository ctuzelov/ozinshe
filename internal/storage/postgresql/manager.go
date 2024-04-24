package storage

import (
	"context"
	"ozinshe/internal/models"
)

type User interface {
	SaveUser(ctx context.Context, user models.User) error
	Delete(ctx context.Context, id int) error
	GetByEmail(ctx context.Context, email string) (models.User, error)
	UpdateTokens(signedToken string, signedRefreshToken string, user_type string) error
	DeleteTokens(ctx context.Context, email string) error
	GetAll(ctx context.Context) ([]models.User, error)
	GetById(ctx context.Context, id int) (models.User, error)
	ChangePassword(ctx context.Context, user models.User) error
	ChangeProfileData(ctx context.Context, user models.User) error
}

type Movie interface {
	GetById(ctx context.Context, id int) (models.Movie, error)
	GetByTitle(ctx context.Context, title string) ([]models.Movie, error)
	GetByYear(ctx context.Context, year_start, year_end int) ([]models.Movie, error)
	GetByGenres(ctx context.Context, genres []string) ([]models.Movie, error)
	GetAll(ctx context.Context) ([]models.Movie, error)
	GetFavorites(ctx context.Context, userID int) ([]models.Movie, error)
	Insert(ctx context.Context, movie models.Movie) (int, error)
	InsertToFavorites(ctx context.Context, movieID, userID int) error
	DeleteFromFavorites(ctx context.Context, movieID, userID int) error
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, movie models.Movie) error
	UpdateCover(ctx context.Context, movieID int, cover models.Cover) error
	UpdateScreenshots(ctx context.Context, movieID int, screenshots []models.Screenshot) error
	FetchGenres(ctx context.Context, movie *models.Movie) error
	FetchMovieData(ctx context.Context) ([]models.Movie, error)
	FetchKeywords(ctx context.Context, movie *models.Movie) error
	FetchCover(ctx context.Context, movie *models.Movie) error
	FetchScreenshots(ctx context.Context, movie *models.Movie) error
	FetchAgeCategories(ctx context.Context, movie *models.Movie) error
}

type Project interface {
	Insert(ctx context.Context, project models.Project) (int, error)
	InsertToFavorites(ctx context.Context, movieID, userID int) error
	DeleteFromFavorites(ctx context.Context, movieID, userID int) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (models.Project, error)
}

type Series interface {
	Insert(ctx context.Context, series models.Series) (int, error)
	InsertToFavorites(ctx context.Context, movieID, userID int) error
	Delete(ctx context.Context, id int) error
	DeleteFromFavorites(ctx context.Context, movieID, userID int) error
	GetById(ctx context.Context, id int) (models.Series, error)
	GetByYear(ctx context.Context, year_start, year_end int) ([]models.Series, error)
	GetByTitle(ctx context.Context, title string) ([]models.Series, error)
	GetByGenres(ctx context.Context, genres []string) ([]models.Series, error)
	GetAll(ctx context.Context) ([]models.Series, error)
	GetFavorites(ctx context.Context, userID int) ([]models.Series, error)
	Update(ctx context.Context, series models.Series) error
	UpdateCover(ctx context.Context, seriesID int, cover models.Cover) error
	UpdateScreenshots(ctx context.Context, seriesID int, screenshots []models.Screenshot) error
	FetchGenres(ctx context.Context, series *models.Series) error
	FetchSeriesData(ctx context.Context) ([]models.Series, error)
	FetchKeywords(ctx context.Context, series *models.Series) error
	FetchCover(ctx context.Context, series *models.Series) error
	FetchScreenshots(ctx context.Context, series *models.Series) error
	FetchAgeCategories(ctx context.Context, series *models.Series) error
	FetchEpisodes(ctx context.Context, seriesID, seasonID int) ([]models.Episode, error)
}

type Genre interface {
	Insert(ctx context.Context, genre []models.Genre) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (models.Genre, error)
	GetAll(ctx context.Context) ([]models.Genre, error)
}

type AgeCategory interface {
	Insert(ctx context.Context, age_category []models.AgeCategory) error
	Delete(ctx context.Context, id int) error
	GetById(ctx context.Context, id int) (models.AgeCategory, error)
	GetAll(ctx context.Context) ([]models.AgeCategory, error)
}

type Keyword interface {
	Insert(ctx context.Context, keywords []models.Keyword) error
}

type Storage struct {
	User
	Movie
	Series
	Genre
	AgeCategory
	Keyword
	Project
}

func New(storage *Postgres) *Storage {
	return &Storage{
		User:        NewUserStorage(storage),
		Genre:       NewGenreStorage(storage),
		Movie:       NewMovieStorage(storage),
		Series:      NewSeriesStorage(storage),
		AgeCategory: NewAgeCategoryStorage(storage),
		Keyword:     NewKeywordStorage(storage),
		Project:     NewProjectStorage(storage),
	}
}
