package handler

import "github.com/gin-gonic/gin"

func (h *Handler) IsAuthorized(c *gin.Context) {
	token, err := c.Cookie("token")

	

	c.Next()
}
