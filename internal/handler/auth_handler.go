package handler

import (
	"errors"
	"net/http"
	"ozinshe/internal/models"
	"ozinshe/internal/validation"

	"github.com/gin-gonic/gin"
)

type entryForm struct {
	Email            string `form:"email"`
	Password         string `form:"password"`
	Confirm_password string `form:"confirm_password"`
}

func (h *Handler) SignUpPage(c *gin.Context) {
	h.render(c, http.StatusOK, "signup.html", nil)
}

func (h *Handler) SignInPage(c *gin.Context) {
	h.render(c, http.StatusOK, "signin.html", nil)
}

func (h *Handler) SignUp(c *gin.Context) {
	var form entryForm

	if err := c.ShouldBind(&form); err != nil {
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

func (h *Handler) SignIn(c *gin.Context) {
	var form entryForm

	if err := c.ShouldBind(&form); err != nil {
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

	c.JSON(http.StatusOK, gin.H{"token": token, "refresh_token": refresh_token})
}

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
