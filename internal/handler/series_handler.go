package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"ozinshe/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Series struct {
	Title       string          `json:"title"`
	Genre       []string        `json:"genres"`
	Year        int             `json:"year"`
	Duration    int             `json:"duration"`
	Keywords    []string        `json:"keywords"` // Consider an array for keywords
	Description string          `json:"description"`
	Producer    string          `json:"producer"`
	Director    string          `json:"director"`
	AgeCategory string          `json:"age_category"`
	Seasons     []models.Season `json:"seasons"`
}

func (h *Handler) CreateSeries(c *gin.Context, form *multipart.Form, project *models.Project) {
	seriesJSON := c.PostForm("series_data")
	var series Series
	err := json.Unmarshal([]byte(seriesJSON), &series)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "error parsing series JSON")
		return
	}

	Genres, AgeCategories, Keywords := ProcessParsing(series.Genre, series.AgeCategory, series.Keywords)

	series_data := models.Series{
		Title:         series.Title,
		ReleaseYear:   series.Year,
		Duration:      series.Duration,
		Genres:        Genres,
		Keywords:      Keywords,
		Description:   series.Description,
		Producer:      series.Producer,
		Director:      series.Director,
		AgeCategories: AgeCategories,
		Seasons:       series.Seasons,
	}

	images_data, err := ProcessSavePhoto(form, "series")
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "processing images in series")
		return
	}

	id, err := h.Service.Series.Add(series_data, images_data)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "adding series")
		return
	}

	project.Project_type = "series"
	project.Project_id = id
}

func (h *Handler) UpdateSeries(c *gin.Context, form *multipart.Form, seriesID int, updated *bool) {
	seriesJSON := c.PostForm("series_data")
	var series Series
	err := json.Unmarshal([]byte(seriesJSON), &series)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "error parsing series JSON")
		return
	}

	Genres, AgeCategories, Keywords := ProcessParsing(series.Genre, series.AgeCategory, series.Keywords)

	series_data := models.Series{
		ID:            seriesID,
		Title:         series.Title,
		ReleaseYear:   series.Year,
		Duration:      series.Duration,
		Genres:        Genres,
		Keywords:      Keywords,
		Description:   series.Description,
		Producer:      series.Producer,
		Director:      series.Director,
		AgeCategories: AgeCategories,
		Seasons:       series.Seasons,
	}

	images_data, err := ProcessSavePhoto(form, "series")
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "processing images in series")
		return
	}

	if form.File["cover"] != nil {
		err := h.Service.Series.UpdateCover(seriesID, images_data)
		if err != nil {
			h.errorpage(c, http.StatusInternalServerError, err, "updating movie cover")
			return
		}
	}

	if form.File["screenshots"] != nil {
		err := h.Service.Series.UpdateScreenshots(seriesID, images_data)
		if err != nil {
			h.errorpage(c, http.StatusInternalServerError, err, "updating movie Screenshots")
			return
		}
	}

	err = h.Service.Series.Update(seriesID, series_data)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "adding series")
		return
	}

	*updated = true

	c.JSON(http.StatusOK, gin.H{"successfully updated with id - ": seriesID})
}

func (h *Handler) GetSeries(c *gin.Context) {
	id := c.Param("seriesID")
	seriesID, _ := strconv.Atoi(id)
	series, err := h.Service.Series.GetById(seriesID)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "series getting failed")
		return
	}
	c.JSON(http.StatusOK, series)
}

func (h *Handler) GetAllSeries(c *gin.Context) {
	series, err := h.Service.Series.GetAll()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "series getting failed")
		return
	}
	c.JSON(http.StatusOK, series)
}
func (h *Handler) GetFilteredSeries(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}

func (h *Handler) DeleteSeries(c *gin.Context) {
	id := c.Param("id")
	seriesID, _ := strconv.Atoi(id)
	err := h.Service.Series.Remove(seriesID)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "deleting series")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Series deleted"})
}

func (h *Handler) GetSeason(c *gin.Context) {
	// 1. Extract Parameters
	seriesID, err := strconv.Atoi(c.Param("seriesID"))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "invalid series id")
		return
	}

	seasonNumber, err := strconv.Atoi(c.Param("seasonNumber"))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "invalid season number")
		return
	}

	// 2. Fetch from Database (Adapt based on your storage logic)
	episodes, err := h.Service.Series.GetSeason(seriesID, seasonNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.errorpage(c, http.StatusNotFound, err, "season not found")
		} else {
			h.errorpage(c, http.StatusInternalServerError, err, "season getting failed")
		}
		return
	}

	c.JSON(http.StatusOK, episodes)
}

func (h *Handler) GetEpisode(c *gin.Context) {
	h.render(c, http.StatusOK, "movie.html", nil)
}
