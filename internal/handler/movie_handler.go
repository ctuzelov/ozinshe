package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"ozinshe/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type movieForm struct {
	Link        string   `json:"movie_link"`
	Title       string   `json:"movie_title"`
	Genres      []string `json:"movie_genres"`
	Year        int      `json:"movie_year"`
	Keywords    []string `json:"movie_keywords"`
	Duration    int      `json:"movie_duration"`
	Producer    string   `json:"movie_producer"`
	Director    string   `json:"movie_director"`
	Description string   `json:"movie_description"`
	AgeCategory string   `json:"movie_age_category"`
}

const maxImageSize = 15 * 1024 * 1024 // 15 MB limit

func (h *Handler) CreateMovie(c *gin.Context, form *multipart.Form, project *models.Project) {
	fmt.Println(c.Request)
	var movie movieForm
	if err := json.NewDecoder(c.Request.Body).Decode(&movie); err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "error parsing movie JSON")
		return
	}

	fmt.Println(movie)
	if movie.Link == "" {
		h.errorpage(c, http.StatusBadRequest, errors.New("empty link"), "empty link")
		return
	}

	Genres, AgeCategories, Keywords := ProcessParsing(movie.Genres, movie.AgeCategory, movie.Keywords)

	movie_data := models.Movie{
		Title:         movie.Title,
		ReleaseYear:   movie.Year,
		Duration:      movie.Duration,
		Description:   movie.Description,
		Producer:      movie.Producer,
		Director:      movie.Director,
		YoutubeID:     movie.Link,
		Genres:        Genres,
		Keywords:      Keywords,
		AgeCategories: AgeCategories,
	}

	images_data, err := ProcessSavePhoto(form, "movies")
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "adding movie")
		return
	}

	id, err := h.Service.Movie.Add(movie_data, images_data)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "adding movie")
		return
	}

	project.Project_type = "movie"
	project.Project_id = id
}

func (h *Handler) UpdateMovie(c *gin.Context, form *multipart.Form, movieID int, updated *bool) {
	movieJSON := c.PostForm("movie_data")
	var movie movieForm
	err := json.Unmarshal([]byte(movieJSON), &movie)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "error parsing series JSON")
		return
	}

	Genres, AgeCategories, Keywords := ProcessParsing(movie.Genres, movie.AgeCategory, movie.Keywords)

	movie_data := models.Movie{
		Title:         movie.Title,
		ReleaseYear:   movie.Year,
		Duration:      movie.Duration,
		Description:   movie.Description,
		Producer:      movie.Producer,
		Director:      movie.Director,
		YoutubeID:     movie.Link,
		Genres:        Genres,
		Keywords:      Keywords,
		AgeCategories: AgeCategories,
	}

	images_data, err := ProcessSavePhoto(form, "movies")
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "adding movie")
		return
	}

	if form.File["cover"] != nil {
		err := h.Service.Movie.UpdateCover(movieID, images_data)
		if err != nil {
			h.errorpage(c, http.StatusInternalServerError, err, "updating movie cover")
			return
		}
	}

	if form.File["screenshots"] != nil {
		err := h.Service.Movie.UpdateScreenshots(movieID, images_data)
		if err != nil {
			h.errorpage(c, http.StatusInternalServerError, err, "updating movie Screenshots")
			return
		}
	}

	err = h.Service.Movie.Update(movieID, movie_data)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "updating movie")
		return
	}

	*updated = true

	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})
}

// @Summary Get details of a specific movie
// @Description Retrieves details of a movie based on the provided ID.
// @Tags movies
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.Movie "Movie details"
// @Failure 400 {object} error "Error getting movie"
// @Router /movies/{id} [get]
func (h *Handler) GetMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	movie, err := h.Service.Movie.GetById(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "movie getting failed")
		return
	}
	c.JSON(http.StatusOK, movie)
}

// @Summary Delete an existing movie
// @Description Deletes an existing movie based on the provided ID. Requires admin authorization.
// @Tags movies
// @Produce json
// @Security CookieAuth
// @Param id path int true "Movie ID"
// @Success 200 "Movie deleted successfully"
// @Failure 500 {object} error "Error deleting movie"
// @Router /movies/{id} [delete]
func (h *Handler) DeleteMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Service.Movie.Remove(id)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "deleting movie")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
}

// @Summary Get a list of all movies
// @Description Retrieves a list of all movies.
// @Tags movies
// @Produce json
// @Success 200 {array} models.Movie "List of movies"
// @Failure 400 {object} error "Error getting movies"
// @Router /movies [get]
func (h *Handler) GetAllMovies(c *gin.Context) {
	movies, err := h.Service.Movie.GetAll()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "movies getting failed")
		return
	}
	c.JSON(http.StatusOK, movies)
}
