package handler

import (
	"net/http"
	"ozinshe/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Get a list of all genres
// @Description Retrieves a list of all genres.
// @Tags genres
// @Produce json
// @Success 200 {array} models.Genre "List of genres"
// @Failure 400 {object} error "Error getting genres"
// @Router /genres [get]
func (h *Handler) GetAllGenres(c *gin.Context) {
	genres, err := h.Service.Genre.GetAll()

	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "genres getting failed")
		return
	}

	c.JSON(http.StatusOK, genres)
}

// @Summary Get details of a specific genre
// @Description Retrieves details of a genre based on the provided ID.
// @Tags genres
// @Produce json
// @Param id path int true "Genre ID"
// @Success 200 {object} models.Genre "Genre details"
// @Failure 400 {object} error "Error getting genre"
// @Router /genres/{id} [get]
func (h *Handler) GetGenre(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	genre, err := h.Service.Genre.GetById(id)

	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "genre getting failed")
		return
	}

	c.JSON(http.StatusOK, genre)
}

// @Summary Create a new genre
// @Description Creates a new genre with the provided name. Requires admin authorization.
// @Tags genres
// @Accept multipart/form-data
// @Produce json
// @Security CookieAuth
// @Param genre formData string true "Genre name"
// @Success 200 "Genre created"
// @Failure 400 {object} error "Error creating genre"
// @Router /genres [post]
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

// @Summary Delete an existing genre
// @Description Deletes an existing genre based on the provided ID. Requires admin authorization.
// @Tags genres
// @Produce json
// @Security CookieAuth
// @Param id path int true "Genre ID"
// @Success 200 "Genre deleted successfully"
// @Failure 400 {object} error "Error deleting genre"
// @Router /genres/{id} [delete]
func (h *Handler) DeleteGenre(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Service.Genre.Remove(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "genre deleting failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Genre deleted"})
}
