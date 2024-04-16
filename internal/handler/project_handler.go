package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func (h *Handler) CreateProject(c *gin.Context) {

	form, err := c.MultipartForm()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding multipart form in create project")
		return
	}

	switch c.PostForm("project_type") {
	case "movie":
		h.CreateMovie(c, form)
	case "series":
		h.CreateSeries(c, form)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project created"})
}

func (h *Handler) GetProject(c *gin.Context) {

}

func (h *Handler) GetFilteredProjects(c *gin.Context) {

}

func (h *Handler) UpdateProject(c *gin.Context) {

}

func (h *Handler) DeleteProject(c *gin.Context) {

}

func (h *Handler) ProjectPage(c *gin.Context) {

}
