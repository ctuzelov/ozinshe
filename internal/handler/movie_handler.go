package handler

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"ozinshe/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type movieForm struct {
	Link        string   `json:"link"`
	Title       string   `json:"title"`
	Genres      []string `json:"genres"`
	Year        int      `json:"year"`
	Keywords    []string `json:"keywords"`
	Duration    int      `json:"duration"`
	Producer    string   `json:"producer"`
	Director    string   `json:"director"`
	Description string   `json:"description"`
	AgeCategory string   `json:"age_category"`
}

const maxImageSize = 15 * 1024 * 1024 // 15 MB limit

func (h *Handler) CreateMovie(c *gin.Context, form *multipart.Form, project *models.Project) {
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

func (h *Handler) GetMovie(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetFilteredMovies(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) DeleteMovie(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Service.Movie.Remove(id)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "deleting movie")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted"})
}

func (h *Handler) GetAllMovies(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}
