package handler

import (
	"errors"
	"net/http"
	"ozinshe/internal/models"
	"ozinshe/internal/validation"

	"github.com/gin-gonic/gin"
)

type entryForm struct {
	Email            string `form:"email" json:"email"`
	Password         string `form:"password" json:"password"`
	Confirm_password string `form:"confirm_password" json:"confirm_password"`
}

// @Summary User sign up
// @Description Registers a new user with the provided email and password.
// @Tags authentication
// @Accept json
// @Produce json
// @Param entry body entryForm true "User email and password"
// @Success 200  "User registered successfully"
// @Failure 400 {object} string "Error signing up: email or password is empty"
// @Failure 400 {object} string "Error signing up: sign up failed"
// @Failure 400 {object} string "Error signing up: user registration failed"
// @Router /signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var form entryForm

	if err := c.ShouldBindJSON(&form); err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding json in sign up")
		return
	}

	user := models.User{
		Email:    form.Email,
		Password: form.Password,
	}

	if user.Email == "" || user.Password == "" {
		h.errorpage(c, http.StatusBadRequest, errors.New("email or password is empty"), "sign up failed")
		return
	}

	err := validation.GetErrMsg(user)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "sign up failed")
		return
	}

	err = h.Service.Register(user)

	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "user registration failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}

// @Summary User sign in
// @Description Logs in an existing user with the provided email and password.
// @Tags authentication
// @Accept json
// @Produce json
// @Param entry body entryForm true "User email and password"
// @Success 200 {array} string "User logged in successfully"
// @Failure 400 {object} string "Error signing in: email or password is empty"
// @Failure 400 {object} string "Error signing in: sign up failed"
// @Failure 400 {object} string "Error signing in: user loggin in failed"
// @Router /signin [post]
func (h *Handler) SignIn(c *gin.Context) {
	var form entryForm

	if err := c.ShouldBindJSON(&form); err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding json in sign in")
		return
	}

	user := models.User{
		Email:    form.Email,
		Password: form.Password,
	}

	if user.Email == "" || user.Password == "" {
		h.errorpage(c, http.StatusBadRequest, errors.New("email or password is empty"), "sign up failed")
		return
	}

	err := validation.GetErrMsg(user)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "sign up failed")
		return
	}

	if form.Password != form.Confirm_password {
		h.errorpage(c, http.StatusBadRequest, errors.New("passwords don't match"), "sign in failed")
		return
	}

	token, refresh_token, err := h.Service.Login(models.User{
		Email:    form.Email,
		Password: form.Password,
	})

	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "user login failed")
		return
	}

	c.SetCookie("token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, []string{token, refresh_token})
}

// @Summary User sign out
// @Description Logs out the currently authenticated user.
// @Tags authentication
// @Security CookieAuth
// @Produce json
// @Success 200 "Logout successful"
// @Failure 400 {object} error "Error logging out user"
// @Router /signout [get]
func (h *Handler) Signout(c *gin.Context) {
	data := c.MustGet("data").(*Data)
	if !data.IsAuthorized {
		return
	}

	err := h.Service.DeleteTokensByEmail(data.User.Email)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "user logout failed")
		return
	}

	c.SetCookie("token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
