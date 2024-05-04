package handler

import (
	"errors"
	"net/http"
	"ozinshe/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type NewPassword struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type ProfileData struct {
	Name        string    `json:"name"`
	Number      string    `json:"number"`
	Password    string    `json:"password"`
	DateOfBirth time.Time `json:"date_of_birth"`
}

// @Summary Get a list of all users
// @Description Retrieves a list of all users in the system. Requires admin authorization.
// @Tags users
// @Security CookieAuth
// @Produce json
// @Success 200 {array} models.User "List of users"
// @Failure 400 {object} error "Error getting users"
// @Router /users [get]
func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.Service.User.GetAll()
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "users getting failed")
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Get details of a specific user
// @Description Retrieves details of a user based on the provided ID.
// @Tags users
// @Security CookieAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User "User details"
// @Failure 400 {object} error "Error getting user"
// @Router /user/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.Service.User.GetById(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "user getting failed")
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Change password
// @Description Allows the currently authenticated user to change their password.
// @Tags users
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param request body NewPassword true "New password information"
// @Success 200 "Password changed successfully"
// @Failure 400 {object} error "Error changing password"
// @Router /change-password [post]
func (h *Handler) ChangePassword(c *gin.Context) {
	data := c.MustGet("data").(*Data)
	if !data.IsAuthorized {
		return
	}

	var form NewPassword
	if err := c.ShouldBindJSON(&form); err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding json in change password")
		return
	}

	err := h.Service.User.UpdatePassword(data.User.Email, form.CurrentPassword, form.NewPassword)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "changing password")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "password changed"})
}

// @Summary Change profile
// @Description Allows the currently authenticated user to change their profile information.
// @Tags users
// @Security CookieAuth
// @Accept json
// @Produce json
// @Param request body ProfileData true "New profile information"
// @Success 200 "Profile changed successfully"
// @Failure 400 {object} error "Error changing profile"
// @Router /change-profile [post]
func (h *Handler) ChangeProfile(c *gin.Context) {
	data := c.MustGet("data").(*Data)
	if !data.IsAuthorized {
		return
	}

	var form ProfileData
	if err := c.ShouldBindJSON(&form); err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "binding json in change profile")
		return
	}

	if form.Name == "" || form.Number == "" || form.DateOfBirth.IsZero() {
		h.errorpage(c, http.StatusBadRequest, errors.New("name, number and date of birth are required"), "changing profile")
		return
	}

	user := models.User{
		Email:       data.User.Email,
		Password:    form.Password,
		Name:        form.Name,
		Number:      form.Number,
		DateOfBirth: form.DateOfBirth,
	}

	err := h.Service.User.UpdateProfile(user)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "changing profile")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile changed"})
}

// @Summary Delete user
// @Description Deletes a user based on the provided ID. Requires admin authorization.
// @Tags users
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Produce json
// @Success 200 "User deleted successfully"
// @Failure 400 {object} error "Error deleting user"
// @Router /user/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	data := c.MustGet("data").(*Data)
	if !data.IsAuthorized {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "invalid user id")
		return
	}

	err = h.Service.User.Remove(id)
	if err != nil {
		h.errorpage(c, http.StatusBadRequest, err, "deleting user")
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
