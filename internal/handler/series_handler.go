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
	Title       string   `json:"series_title"`
	Genre       []string `json:"series_genres"`
	Year        int      `json:"series_year"`
	Duration    int      `json:"series_duration"`
	Keywords    []string `json:"series_keywords"` // Consider an array for keywords
	Description string   `json:"series_description"`
	Producer    string   `json:"series_producer"`
	Director    string   `json:"series_director"`
	AgeCategory string   `json:"series_age_category"`
	Seasons     []Season `json:"series_seasons"`
}

type Season struct {
	Episode []Episode `json:"season_episodes"`
}

type Episode struct {
	Link string `json:"episode_link"`
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
	Seasons := ProcessParsingSeasons(series.Seasons)

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
		Seasons:       Seasons,
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
	Seasons := ProcessParsingSeasons(series.Seasons)

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
		Seasons:       Seasons,
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

// @Summary Get details of a specific series
// @Description Retrieves details of a series based on the provided ID.
// @Tags series
// @Produce json
// @Param seriesID path string true "Series ID"
// @Success 200 {object} models.Series "Series details"
// @Failure 400 {object} error"Error getting series"
// @Router /series/{seriesID} [get]
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

// @Summary Get a list of all series
// @Description Retrieves a list of all series.
// @Tags series
// @Produce json
// @Success 200 {array} models.Series "List of series"
// @Failure 400 {object} error "Error getting series"
// @Router /series [get]
func (h *Handler) GetAllSeries(c *gin.Context) {
	series, err := h.Service.Series.GetAll()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "series getting failed")
		return
	}
	c.JSON(http.StatusOK, series)
}

// @Summary Delete a series
// @Description Deletes a series based on the provided ID.
// @Tags series
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Series ID"
// @Success 200 "Series deleted successfully"
// @Failure 400 {object} error "Invalid parameters"
// @Failure 500 {object} error "Error deleting series"
// @Router /series/{id} [delete]
func (h *Handler) DeleteSeries(c *gin.Context) {
	id := c.Param("id")
	seriesID, err := strconv.Atoi(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "invalid parameters")
		return
	}
	err = h.Service.Series.Remove(seriesID)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "deleting series")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Series deleted"})
}

// @Summary Get episodes of a specific season of a series
// @Description Retrieves episodes of a specific season of a series based on the provided series ID and season number.
// @Tags series
// @Produce json
// @Param seriesID path string true "Series ID"
// @Param seasonNumber path string true "Season Number"
// @Success 200 {array} models.Episode "List of episodes"
// @Failure 400 {object} error "Invalid parameters"
// @Failure 404 {object} error "Season not found"
// @Failure 500 {object} error "Error getting episodes"
// @Router /series/{seriesID}/seasons/{seasonNumber}/episodes [get]
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

// @Summary Get details of a specific episode of a series
// @Description Retrieves details of a specific episode of a series based on the provided series ID, season number, and episode ID.
// @Tags series
// @Produce json
// @Param seriesID path string true "Series ID"
// @Param seasonNumber path string true "Season Number"
// @Param episodeID path string true "Episode ID"
// @Success 200 {object} models.Episode "Episode details"
// @Failure 400 {object} error "Invalid parameters"
// @Failure 500 {object} error "Error getting episode"
// @Router /series/{seriesID}/seasons/{seasonNumber}/episodes/{episodeID} [get]
func (h *Handler) GetEpisode(c *gin.Context) {
	var err error
	seriesID, err := strconv.Atoi(c.Param("seriesID"))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "invalid parameters")
		return
	}

	seasonNumber, err := strconv.Atoi(c.Param("seasonNumber"))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "invalid parameters")
		return
	}
	episodeID, err := strconv.Atoi(c.Param("episodeID"))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "invalid parameters")
		return
	}

	episode, err := h.Service.Series.GetEpisode(seriesID, seasonNumber, episodeID)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "episode getting failed")
	}

	c.JSON(http.StatusOK, episode)
}
