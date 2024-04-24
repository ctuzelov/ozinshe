package handler

import (
	"net/http"
	"ozinshe/internal/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func (h *Handler) GetAgeCategory(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	ageCategory, err := h.Service.AgeCategory.GetById(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "age category getting failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Age category gotten  ", "age_category": ageCategory})
}

func (h *Handler) DeleteAgeCategory(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))
	err := h.Service.AgeCategory.Remove(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "age category deleting failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Age category deleted"})
}

func (h *Handler) GetAllAgeCategories(c *gin.Context) {

	ageCategories, err := h.Service.AgeCategory.GetAll()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "age categories getting failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Age categories gotten", "age_categories": ageCategories})
}
