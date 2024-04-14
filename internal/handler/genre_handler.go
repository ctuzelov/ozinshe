package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GenresPage(c *gin.Context) {
	h.render(c, http.StatusOK, "genres.html", nil)
}

func (h *Handler) GetAllGenres(c *gin.Context) {
	h.render(c, http.StatusOK, "genres.html", nil)
}

func (h *Handler) GetGenre(c *gin.Context) {
	h.render(c, http.StatusOK, "genres.html", nil)
}

func (h *Handler) CreateGenre(c *gin.Context) {
	h.render(c, http.StatusOK, "genres.html", nil)
}

func (h *Handler) DeleteGenre(c *gin.Context) {
	h.render(c, http.StatusOK, "genres.html", nil)
}
