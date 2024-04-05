package handler

import (
	"errors"
	"net/http"
	"ozinshe/internal/models"

	"github.com/gin-gonic/gin"
)

type entryForm struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	Confirm_password string `json:"confirm_password"`
}

func (h *Handler) SignUp(c *gin.Context) {
	var form entryForm

	if err := c.ShouldBindJSON(&form); err != nil {
		// TODO: send status code
		h.Log.Error("error in binding json", err)
		return
	}

	// TODO: validate form

	err := h.Service.Register(models.User{
		Email:    form.Email,
		Password: form.Password,
	})

	if err != nil {
		// TODO: send status code
		h.Log.Error("User registration failed", err)
	}
}

func (h *Handler) SignIn(c *gin.Context) {
	var form entryForm

	if err := c.ShouldBindJSON(&form); err != nil {
		// TODO: send status code
		h.Log.Error("error in binding json", err)
		return
	}

	// TODO: validate form

	if form.Password != form.Confirm_password {
		// TODO: send status code
		h.Log.Error("User login failed", errors.New("passwords don't match"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords don't match"})
		return
	}

	token, refresh_token, err := h.Service.Login(models.User{
		Email:    form.Email,
		Password: form.Password,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.Log.With("User login failed", err)
		return
	}

	c.SetCookie("token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refresh_token})
}

func (h *Handler) Signout(c *gin.Context) {
    // Clear token cookie
    c.SetCookie("token", "", -1, "/", "", false, true) 

    // Clear refresh token cookie
    c.SetCookie("refresh_token", "", -1, "/", "", false, true) 

    c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}