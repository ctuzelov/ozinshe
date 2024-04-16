package handler

import (
	"mime/multipart"
	"net/http"
	"os"
	"ozinshe/internal/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type movieForm struct {
	Link        string `form:"link"`
	Title       string `form:"title"`
	Genres      string `form:"genre"`
	Year        int    `form:"year"`
	Keywords    string `form:"keywords"`
	Duration    int    `form:"duration"`
	Producer    string `form:"producer"`
	Director    string `form:"director"`
	Description string `form:"description"`
	AgeCategory string `form:"age_category"`
}

const maxImageSize = 15 * 1024 * 1024 // 2 MB limit

func (h *Handler) CreateMovie(c *gin.Context, form *multipart.Form) {
	var movie_form movieForm

	if err := c.ShouldBind(&movie_form); err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding form in create movie")
		return
	}

	movie_genres, movie_ages, movie_keywords := ProcessParsing(movie_form.Genres, movie_form.Keywords, movie_form.AgeCategory)

	movie := models.Movie{
		Title:         movie_form.Title,
		ReleaseYear:   movie_form.Year,
		Duration:      movie_form.Duration,
		Genres:        movie_genres,
		Keywords:      movie_keywords,
		Description:   movie_form.Description,
		Producer:      movie_form.Producer,
		Director:      movie_form.Director,
		YoutubeID:     movie_form.Link,
		AgeCategories: movie_ages,
	}

	currentDir, err := os.Getwd()
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "getting current dir")
		return
	}
	path := currentDir + "/../uploads/movies/"

	images_data := models.SavePhoto{
		File_form:    form,
		UploadPath:   path,
		MaxImageSize: maxImageSize,
	}

	id, err := h.Service.Movie.Add(movie, images_data)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "adding movie")
		return
	}

	c.JSON(http.StatusOK, gin.H{"successfully created with id - ": id})
}

func ProcessParsing(Genres, AgeCategory, Keywords string) ([]models.Genre, []models.AgeCategory, []models.Keyword) {

	genres := strings.Split(Genres, ",")
	keywords := strings.Split(Keywords, ",")
	ageCategory := strings.Split(AgeCategory, ",")

	movie_genres := []models.Genre{}
	for _, genre := range genres {
		movie_genres = append(movie_genres, models.Genre{Name: genre})
	}

	movie_ages := []models.AgeCategory{}
	for _, age := range ageCategory {
		ages := strings.Split(age, "-")
		min_age, _ := strconv.Atoi(ages[0])
		max_age, _ := strconv.Atoi(ages[1])
		movie_ages = append(movie_ages, models.AgeCategory{MinAge: min_age, MaxAge: max_age})
	}

	movie_keywords := []models.Keyword{}
	for _, keyword := range keywords {
		movie_keywords = append(movie_keywords, models.Keyword{Name: keyword})
	}

	return movie_genres, movie_ages, movie_keywords
}

func (h *Handler) MoviePage(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetMovie(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetFilteredMovies(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) UpdateMovie(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) DeleteMovie(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetAllMovies(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}
