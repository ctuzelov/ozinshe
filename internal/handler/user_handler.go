package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetAllUsers(c *gin.Context) {
	h.render(c, http.StatusOK, "users.html", nil)
}

func (h *Handler) GetUser(c *gin.Context) {
	h.render(c, http.StatusOK, "users.html", nil)
}
