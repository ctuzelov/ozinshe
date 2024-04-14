package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) SeriesPage(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetSeries(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetFilteredSeries(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) CreateSeries(c *gin.Context) {
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
