package handler

import (
	"errors"
	"net/http"
	"ozinshe/internal/models"
	"ozinshe/util"
	"strconv"

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

		id, _ := strconv.Atoi(validToken["uid"].(string))
		data.User.Email, data.User.UserType, data.User.Name, data.User.ID = validToken["email"].(string), validToken["user_type"].(string), validToken["name"].(string), id
		data.IsAuthorized = true
		data.IsAdmin = data.User.UserType == "admin"
	}

	c.Set("data", data)

	c.Next()
}

func (h *Handler) IsAdminMiddlware(c *gin.Context) {
	data := c.MustGet("data").(*Data)

	if !data.IsAdmin {
		h.errorpage(c, http.StatusForbidden, errors.New("admin access required"), "forbidden")
		c.Abort()
		return
	}

	c.Next()
}

func (h *Handler) MustBeAuthorizedMiddleware(c *gin.Context) {
	data := c.MustGet("data").(*Data)

	if !data.IsAuthorized {
		h.errorpage(c, http.StatusForbidden, nil, "forbidden")
		return
	}

	c.Next()
}
