package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"ozinshe/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service   *service.Service
	Log       *slog.Logger
	TempCache *template.Template
}

func New(service *service.Service, log *slog.Logger) *Handler {
	cache, _ := template.ParseGlob("ui/html/*.html") // TODO: handle the potential error
	return &Handler{
		Service:   service,
		Log:       log,
		TempCache: cache,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.StaticFS("/ui/assets/", http.Dir("./ui/assets/"))

	router.Use(h.Middleware)

	router.GET("/", h.HomePage)
	router.GET("/signup", h.SignUpPage)
	router.GET("/signin", h.SignInPage)
	router.POST("/signup", h.SignUp)
	router.POST("/signin", h.SignIn)

	router.GET("/signout", h.Signout)

	return router
}
