package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) MoviePage(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetMovie(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) GetFilteredMovies(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) CreateMovie(c *gin.Context) {
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

