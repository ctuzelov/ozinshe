package handler

import (
	"net/http"
	"ozinshe/internal/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// @Summary Create a new age category
// @Description Creates a new age category based on the provided age range. Requires admin authorization.
// @Tags age categories
// @Accept multipart/form-data
// @Produce json
// @Security CookieAuth
// @Param age formData string true "Age range (e.g., '18-30')"
// @Success 200 "Age category created"
// @Failure 400 {object} error "Error creating age category"
// @Router /ages [post]
func (h *Handler) CreateAgeCategory(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding multipart form in create age category")
		return
	}

	ages := form.Value["age"]
	minAgeString, maxAgeString := strings.Split(ages[0], "-")[0], strings.Split(ages[0], "-")[1]

	minAge, _ := strconv.Atoi(minAgeString)
	maxAge, _ := strconv.Atoi(maxAgeString)

	ageCategory := models.AgeCategory{
		MinAge: minAge,
		MaxAge: maxAge,
	}

	err = h.Service.AgeCategory.Add(append([]models.AgeCategory{}, ageCategory))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "age category creation failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Age category created"})
}

// @Summary Get details of a specific age category
// @Description Retrieves details of an age category based on the provided ID.
// @Tags age categories
// @Produce json
// @Param id path int true "Age category ID"
// @Success 200 {object} models.AgeCategory "Age category retrieved successfully"
// @Failure 400 {object} error "Error retrieving age category"
// @Router /ages/{id} [get]
func (h *Handler) GetAgeCategory(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	ageCategory, err := h.Service.AgeCategory.GetById(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "age category getting failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Age category gotten  ", "age_category": ageCategory})
}

// @Summary Delete an existing age category
// @Description Deletes an existing age category based on the provided ID.Requires admin authorization.
// @Tags age categories
// @Produce json
// @Security CookieAuth
// @Param id path int true "Age category ID"
// @Success 200 "Age category deleted successfully"
// @Failure 400 {object} error "Error deleting age category"
// @Router /ages/{id} [delete]
func (h *Handler) DeleteAgeCategory(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Service.AgeCategory.Remove(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "age category deleting failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Age category deleted"})
}

// @Summary Get a list of all age categories
// @Description Retrieves a list of all age categories in the system.
// @Tags age categories
// @Produce json
// @Success 200 {array} models.AgeCategory "List of age categories"
// @Failure 400 {object} error "Error retrieving age categories"
// @Router /ages [get]
func (h *Handler) GetAllAgeCategories(c *gin.Context) {

	ageCategories, err := h.Service.AgeCategory.GetAll()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "age categories getting failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Age categories gotten", "age_categories": ageCategories})
}
