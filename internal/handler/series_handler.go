package handler

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Season struct {
	SeasonNumber int       `json:"season_number"`
	Episodes     []Episode `json:"episodes"`
}

type Episode struct {
	EpisodeNumber int    `json:"episode_number"`
	Link          string `json:"link"`
}

type Series struct {
	Title       string   `json:"title"`
	Genre       string   `json:"genre"`
	Year        int      `json:"year"`
	Duration    int      `json:"duration"`
	Keywords    []string `json:"keywords"` // Consider an array for keywords
	Description string   `json:"description"`
	Producer    string   `json:"producer"`
	Director    string   `json:"director"`
	AgeCategory string   `json:"age_category"`
	Seasons     []Season `json:"seasons"`
}

func (h *Handler) SeriesPage(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetSeries(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetFilteredSeries(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) CreateSeries(c *gin.Context, form *multipart.Form) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) UpdateSeries(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) DeleteSeries(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetAllSeries(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetSeason(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetEpisode(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}
