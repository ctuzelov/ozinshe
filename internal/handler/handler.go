package handler

import (
	"log/slog"
	"ozinshe/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
	Log     *slog.Logger
}

func New(service *service.Service, log *slog.Logger) *Handler {
	return &Handler{
		Service: service,
		Log:     log,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/signup", h.SignUp)
	router.POST("/signin", h.SignIn)

	return router
}
