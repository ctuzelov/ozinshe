package handler

import (
	"net/http"
	"ozinshe/internal/models"
	"ozinshe/util"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Middleware(c *gin.Context) {
	token, err := c.Cookie("token")
	data := &Data{}

	switch err {
	case http.ErrNoCookie:
		data.User = models.User{}
		data.IsAuthorized = false
		data.IsAdmin = false
	case nil:
		validToken, err := util.ValidateToken(token)
		if err != nil {
			h.errorpage(c, http.StatusBadRequest, err, "token validation failed")
		}
		data.User.Email, data.User.UserType, data.User.Name = validToken["email"].(string), validToken["user_type"].(string), validToken["name"].(string)
		data.IsAuthorized = true
		data.IsAdmin = data.User.UserType == "admin"
	}

	c.Set("data", data)

	c.Next()
}
