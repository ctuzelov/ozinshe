package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"ozinshe/internal/models"

	"github.com/gin-gonic/gin"
)

type Data struct {
	User         models.User
	Content      any
	IsAuthorized bool
	IsAdmin      bool
	ErrMsgs      map[string]string
}

type ErrorData struct {
	Status  int
	Message string
}

func (h *Handler) render(c *gin.Context, status int, page string, data any) {
	buf := new(bytes.Buffer)

	err := h.TempCache.ExecuteTemplate(buf, page, data)
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, fmt.Sprintf("template error: %s", page)) // Log template name
		return
	}

	// Handle potential error writing the response
	_, err = c.Writer.Write(buf.Bytes())
	if err != nil {
		h.errorpage(c, http.StatusInternalServerError, err, "error writing response")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8") // Set header explicitly
	c.Status(status)                                     // Set status code
}

func (h *Handler) errorpage(c *gin.Context, status int, err error, errortype string) {
	if err != nil {
		// Consider whether to unwrap or keep error wrapping depending on your debugging needs
		h.Log.Error("%s: %v", errortype, err) // Clearer formatting
	}

	// Customize error data based on the situation
	errdata := ErrorData{
		Status:  status,
		Message: http.StatusText(status), // Default message
	}

	// Optionally provide a more user-friendly message
	if status == http.StatusInternalServerError {
		errdata.Message = "Something went wrong. Please try again later."
	}

	h.render(c, status, "error.html", errdata)
}
