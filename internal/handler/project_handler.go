package handler

import (
	"fmt"
	"net/http"
	"ozinshe/internal/models"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Contents struct {
	Movies []models.Movie
	Series []models.Series
}

func (h *Handler) CreateProject(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding multipart form in create project")
		return
	}

	project := models.Project{}
	switch c.PostForm("project_type") {
	case "movie":
		h.CreateMovie(c, form, &project)
	case "series":
		h.CreateSeries(c, form, &project)
	}

	if project.Project_type == "" {
		return
	}

	id, err := h.Service.Project.Add(project)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "project creation failed")
		return
	}

	if project.Project_type == "movie" {
		movie_data, err := h.Service.Movie.GetById(project.Project_id)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "movie getting failed")
			return
		}
		movie_data.ID = id
		c.JSON(http.StatusOK, movie_data)
	} else if project.Project_type == "series" {
		series_data, err := h.Service.Series.GetById(project.Project_id)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "series getting failed")
			return
		}
		series_data.ID = id
		c.JSON(http.StatusOK, series_data)
	}
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "invalid project id")
		return
	}

	project, err := h.Service.Project.GetById(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "project getting failed")
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding multipart form in update project")
		return
	}

	updated := false
	switch project.Project_type {
	case "movie":
		h.UpdateMovie(c, form, project.Project_id, &updated)
	case "series":
		h.UpdateSeries(c, form, project.Project_id, &updated)
	}

	if !updated {
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated"})
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	project, err := h.Service.Project.GetById(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "project getting failed")
		return
	}

	if project.Project_type == "movie" {
		err := h.Service.Movie.Remove(project.Project_id)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "movie removing failed")
			return
		}
	} else if project.Project_type == "series" {
		err := h.Service.Series.Remove(project.Project_id)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "series removing failed")
			return
		}
	}

	err = h.Service.Project.Remove(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "project deleting failed")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted"})
}

func (h *Handler) GetProject(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	project, err := h.Service.Project.GetById(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "project getting failed")
		return
	}

	switch project.Project_type {
	case "movie":
		movie, err := h.Service.Movie.GetById(project.Project_id)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "movie getting failed")
			return
		}
		c.JSON(http.StatusOK, movie)

	case "series":
		series, err := h.Service.Series.GetById(project.Project_id)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "series getting failed")
			return
		}
		c.JSON(http.StatusOK, series)
	}
}

func (h *Handler) GetAllProjects(c *gin.Context) {
	movies, err := h.Service.Movie.GetAll()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "(movie)projects getting failed")
		return
	}
	var contents Contents

	series, err := h.Service.Series.GetAll()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "(series)projects getting failed")
		return
	}

	contents.Series = series
	contents.Movies = movies

	c.JSON(http.StatusOK, contents)
}

func (h *Handler) GetAllFavorites(c *gin.Context) {
	data := c.MustGet("data").(*Data)

	var content Contents
	movies, err := h.Service.Movie.GetFavorites(data.User.ID)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "getting favorites failed")
		return
	}
	content.Movies = movies

	series, err := h.Service.Series.GetFavorites(data.User.ID)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "getting favorites failed")
		return
	}
	content.Series = series

	c.JSON(http.StatusOK, content)
}

func (h *Handler) GetFilteredProjects(c *gin.Context) {
	var filter models.FilterParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding query failed")
		return
	}

	var contents Contents

	if filter.YearEnd < filter.YearStart || filter.YearEnd == 0 {
		filter.YearEnd = time.Now().Year()
	}

	fmt.Println(filter)

	switch filter.Project_type {
	case "movie":
		movies, err := h.Service.Movie.GetFiltered(filter)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "movie filtering failed")
			return
		}
		if filter.PopularityOrder == "asc" {
			sort.SliceStable(movies, func(i, j int) bool {
				return movies[i].Popularity > movies[j].Popularity
			})
		} else if filter.PopularityOrder == "desc" {
			sort.SliceStable(movies, func(i, j int) bool {
				return movies[i].Popularity < movies[j].Popularity
			})
		}

		c.JSON(http.StatusOK, movies)

	case "series":
		series, err := h.Service.Series.GetFiltered(filter)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "series filtering failed")
			return
		}

		if filter.PopularityOrder == "asc" {
			sort.SliceStable(series, func(i, j int) bool {
				return series[i].Popularity > series[j].Popularity
			})
		} else if filter.PopularityOrder == "desc" {
			sort.SliceStable(series, func(i, j int) bool {
				return series[i].Popularity < series[j].Popularity
			})
		}
		c.JSON(http.StatusOK, series)

	default:
		movies, err := h.Service.Movie.GetFiltered(filter)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "movie filtering failed")
			return
		}

		series, err := h.Service.Series.GetFiltered(filter)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "series filtering failed")
			return
		}

		contents.Movies = movies
		contents.Series = series

		if filter.PopularityOrder == "asc" {
			sort.SliceStable(contents.Movies, func(i, j int) bool {
				return contents.Movies[i].Popularity > contents.Movies[j].Popularity
			})
		} else if filter.PopularityOrder == "desc" {
			sort.SliceStable(contents.Movies, func(i, j int) bool {
				return contents.Movies[i].Popularity < contents.Movies[j].Popularity
			})
		}

		if filter.PopularityOrder == "asc" {
			sort.SliceStable(contents.Series, func(i, j int) bool {
				return contents.Series[i].Popularity > contents.Series[j].Popularity
			})
		} else if filter.PopularityOrder == "desc" {
			sort.SliceStable(contents.Series, func(i, j int) bool {
				return contents.Series[i].Popularity < contents.Series[j].Popularity
			})
		}

		c.JSON(http.StatusOK, contents)
	}

}

func (h *Handler) MakeFavorite(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	data := c.MustGet("data").(*Data)

	project, err := h.Service.Project.GetById(projectID)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "project getting failed")
		return
	}

	if project.Project_type == "movie" {
		err := h.Service.Movie.AddToFavorites(project.Project_id, data.User.ID)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "adding to favorites failed")
			return
		}
	} else if project.Project_type == "series" {
		err := h.Service.Series.AddToFavorites(project.Project_id, data.User.ID)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "adding to favorites failed")
			return
		}
	}

	err = h.Service.Project.AddToFavorites(projectID, data.User.ID)

	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "adding to favorites failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "added to favorites"})
}

func (h *Handler) RemoveFromFavorites(c *gin.Context) {
	projectID, _ := strconv.Atoi(c.Param("id"))
	data := c.MustGet("data").(*Data)

	project, err := h.Service.Project.GetById(projectID)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "project getting failed")
		return
	}

	if project.Project_type == "movie" {
		err := h.Service.Movie.RemoveFromFavorites(project.Project_id, data.User.ID)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "removing from favorites failed")
			return
		}
	} else if project.Project_type == "series" {
		err := h.Service.Series.RemoveFromFavorites(project.Project_id, data.User.ID)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "removing from favorites failed")
			return
		}
	}

	err = h.Service.Project.RemoveFromFavorites(projectID, data.User.ID)

	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "removing from favorites failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "removed from favorites"})
}
