package handler

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"ozinshe/internal/models"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Data struct {
	User         models.User
	Content      any
	IsAuthorized bool
	IsAdmin      bool
	ErrMsgs      map[string]string
}

type ErrorData struct {
	Status  int
	Message string
}

func (h *Handler) render(c *gin.Context, status int, page string, data any) {
	buf := new(bytes.Buffer)

	err := h.TempCache.ExecuteTemplate(buf, page, data)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, fmt.Sprintf("template error: %s", page))
		return
	}

	_, err = c.Writer.Write(buf.Bytes())
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "error writing response")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(status)
}

func (h *Handler) errorpage(c *gin.Context, status int, err error, errortype string) {
	if err != nil {
		h.Log.Error(fmt.Sprintf("%s: %v", errortype, err))
	}

	errdata := ErrorData{
		Status:  status,
		Message: getLastMessageFromError(err),
	}

	if status == http.StatusInternalServerError {
		errdata.Message = "Something went wrong. Please try again later."
	}

	c.JSON(status, errdata)
	c.Abort()
}

func getLastMessageFromError(err error) string {
	if err == nil {
		return ""
	}
	// Convert the error to a string
	errStr := err.Error()

	// Define a regular expression pattern to match the last message after the last colon
	re := regexp.MustCompile(`[^:]+:\s*([^:]+)$`)

	// Find the last match of the pattern in the error string
	matches := re.FindStringSubmatch(errStr)

	// If there is a match, return the captured group (message)
	if len(matches) > 1 {
		return matches[1]
	}

	// If no match found, return the entire error string
	return errStr
}

func ProcessSavePhoto(form *multipart.Form, dst string) (models.SavePhoto, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return models.SavePhoto{}, err
	}
	path := currentDir + "/internal/uploads/" + dst + "/"

	images_data := models.SavePhoto{
		File_form:    form,
		UploadPath:   path,
		MaxImageSize: maxImageSize,
	}

	return images_data, nil
}

func ProcessParsing(genres []string, ageCategory string, keywords []string) ([]models.Genre, []models.AgeCategory, []models.Keyword) {
	var Keywords []models.Keyword
	for _, keyword := range keywords {
		keyword := models.Keyword{
			Name: keyword,
		}
		Keywords = append(Keywords, keyword)
	}

	var AgeCategories []models.AgeCategory
	ages := strings.Split(ageCategory, "-")
	min_age, _ := strconv.Atoi(ages[0])
	max_age, _ := strconv.Atoi(ages[1])
	age_category := models.AgeCategory{
		MinAge: min_age,
		MaxAge: max_age,
	}
	AgeCategories = append(AgeCategories, age_category)

	var Genres []models.Genre
	for _, genre := range genres {
		genre := models.Genre{
			Name: genre,
		}
		Genres = append(Genres, genre)
	}

	return Genres, AgeCategories, Keywords
}

func ProcessParsingSeasons(seasons []Season) []models.Season {
	var Seasons []models.Season

	for i, season := range seasons {
		var Episodes []models.Episode
		for i, episode := range season.Episode {
			episode := models.Episode{
				EpisodeNumber: i + 1,
				Link:          episode.Link,
			}
			Episodes = append(Episodes, episode)
		}

		season := models.Season{
			SeasonNumber: i + 1,
			Episodes:     Episodes,
		}

		Seasons = append(Seasons, season)

	}

	return Seasons
}
