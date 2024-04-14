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

	router.GET("/", h.HomePage)
	router.GET("/signup", h.SignUpPage)
	router.GET("/signin", h.SignInPage)
	router.POST("/signup", h.SignUp)
	router.POST("/signin", h.SignIn)

	authGroup := router.Group("/", h.Middleware)
	{
		authGroup.GET("/signout", h.Signout)

		authGroup.GET("/users", h.GetAllUsers)
		authGroup.GET("/user/:id", h.GetUser)

		projectGroup := authGroup.Group("/projects")
		{
			projectGroup.GET("/", h.GetFilteredProjects)
			projectGroup.POST("/create-project", h.CreateProject)
			projectGroup.GET("/:id", h.ProjectPage)
			projectGroup.DELETE("/:id", h.DeleteProject)
			projectGroup.PUT("/:id", h.UpdateProject)
		}

		movieGroup := authGroup.Group("/movies")
		{
			movieGroup.GET("/", h.GetFilteredMovies)
			movieGroup.GET("/:id", h.GetMovie)
			movieGroup.POST("/create-movie", h.CreateMovie)
			movieGroup.DELETE("/:id", h.DeleteMovie)
			movieGroup.PUT("/:id", h.UpdateMovie)
		}

		seriesGroup := authGroup.Group("/series")
		{
			seriesGroup.GET("/", h.GetFilteredSeries)
			seriesGroup.GET("/:id", h.GetSeries)
			seriesGroup.POST("/create-series", h.CreateSeries)
			seriesGroup.DELETE("/:id", h.DeleteSeries)
			seriesGroup.PUT("/:id", h.UpdateSeries)
			
			seasonGroup := seriesGroup.Group("/seasons")
			{
				seasonGroup.GET("/:id", h.GetSeason)
			}

			episodeGroup := seriesGroup.Group("/episodes")
			{
				episodeGroup.GET("/:id", h.GetEpisode)
			}
		}

		genreGroup := authGroup.Group("/genres")
		{
			genreGroup.GET("/", h.GetAllGenres)
			genreGroup.GET("/:id", h.GetGenre)
			genreGroup.POST("/", h.CreateGenre)
			genreGroup.DELETE("/:id", h.DeleteGenre)
		}
	}

	return router
}
