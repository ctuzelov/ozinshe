package handler

import (
	"net/http"
	"ozinshe/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GenresPage(c *gin.Context) {
	h.render(c, http.StatusOK, "genres.html", nil)
}

func (h *Handler) GetAllGenres(c *gin.Context) {
	genres, err := h.Service.Genre.GetAll()

	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "genres getting failed")
		return
	}

	c.JSON(http.StatusOK, genres)
}

func (h *Handler) GetGenre(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	genre, err := h.Service.Genre.GetById(id)

	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "genre getting failed")
		return
	}

	c.JSON(http.StatusOK, genre)
}

func (h *Handler) CreateGenre(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding multipart form in create genre")
		return
	}

	genres := form.Value["genre"]

	err = h.Service.Genre.Add(append([]models.Genre{}, models.Genre{Name: genres[0]}))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "genre creation failed")
		return
	}
}

func (h *Handler) DeleteGenre(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Service.Genre.Remove(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "genre deleting failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Genre deleted"})
}
