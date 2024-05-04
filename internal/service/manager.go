package service

import (
	"ozinshe/internal/models"
	psql "ozinshe/internal/storage/postgresql"
)

type User interface {
	Register(user models.User) error
	Login(user models.User) (string, string, error)
	UpdateAllTokens(signedToken string, signedRefreshToken string, user_type string, id string) (string, string, error)
	Remove(id int) error
	DeleteTokensByEmail(email string) error
	UpdatePassword(email string, current_password, new_password string) error
	UpdateProfile(user models.User) error
	GetById(id string) (models.User, error)
	GetAll() ([]models.User, error)
}

type Movie interface {
	Add(movie models.Movie, image_data models.SavePhoto) (int, error)
	AddToFavorites(movieID, userID int) error
	RemoveFromFavorites(movieID, userID int) error
	Remove(id int) error
	Update(id int, movie models.Movie) error
	UpdateCover(id int, image_data models.SavePhoto) error
	UpdateScreenshots(id int, image_data models.SavePhoto) error
	GetById(id int) (models.Movie, error)
	GetAll() ([]models.Movie, error)
	GetFavorites(userID int) ([]models.Movie, error)
	GetFiltered(filter models.FilterParams) ([]models.Movie, error)
}

type Series interface {
	Add(series models.Series, image_data models.SavePhoto) (int, error)
	AddToFavorites(movieID, userID int) error
	RemoveFromFavorites(movieID, userID int) error
	Remove(id int) error
	GetById(id int) (models.Series, error)
	GetAll() ([]models.Series, error)
	GetFavorites(userID int) ([]models.Series, error)
	GetSeason(seriesID, seasonNumber int) ([]models.Episode, error)
	GetEpisode(seriesID, seasonNumber, episodeNumber int) (models.Episode, error)
	Update(id int, series models.Series) error
	UpdateCover(id int, image_data models.SavePhoto) error
	UpdateScreenshots(id int, image_data models.SavePhoto) error
	GetFiltered(filter models.FilterParams) ([]models.Series, error)
}

type Genre interface {
	Add(genres []models.Genre) error
	Remove(id int) error
	GetById(id int) (models.Genre, error)
	GetAll() ([]models.Genre, error)
}

type AgeCategory interface {
	Add(age_category []models.AgeCategory) error
	Remove(id int) error
	GetById(id int) (models.AgeCategory, error)
	GetAll() ([]models.AgeCategory, error)
}

type Keyword interface {
	Add(keywords []models.Keyword) error
}

type Project interface {
	Add(project models.Project) (int, error)
	Remove(id int) error
	GetById(id int) (models.Project, error)
	AddToFavorites(movieID, userID int) error
	RemoveFromFavorites(movieID, userID int) error
}

type Service struct {
	User
	Movie
	Series
	Genre
	Keyword
	AgeCategory
	Project
}

func New(storage *psql.Storage) *Service {
	return &Service{
		User:        NewUserService(storage.User),
		Movie:       NewMovieService(storage.Movie),
		Series:      NewSeriesService(storage.Series),
		Genre:       NewGenreService(storage.Genre),
		Keyword:     NewKeywordService(storage.Keyword),
		AgeCategory: NewAgeCategoryService(storage.AgeCategory),
		Project:     NewProjectService(storage.Project),
	}
}
