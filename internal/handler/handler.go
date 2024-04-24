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
		authGroup.GET("/signout", h.MustBeAuthorizedMiddleware, h.Signout)

		authGroup.GET("/users", h.GetAllUsers)
		authGroup.GET("/user/:id", h.GetUser)
		authGroup.POST("/change-password", h.MustBeAuthorizedMiddleware, h.ChangePassword)
		authGroup.POST("/change-profile", h.MustBeAuthorizedMiddleware, h.ChangeProfile)
		authGroup.DELETE("/user/:id", h.IsAdminMiddlware, h.DeleteUser)

		projectGroup := authGroup.Group("/projects")
		{
			projectGroup.GET("/", h.GetAllProjects)
			projectGroup.GET("/:id", h.GetProject)
			projectGroup.GET("/search", h.GetFilteredProjects)
			projectGroup.GET("/favorites", h.MustBeAuthorizedMiddleware, h.GetAllFavorites)
			projectGroup.POST("/:id/favorite", h.MustBeAuthorizedMiddleware, h.MakeFavorite)
			projectGroup.DELETE("/:id/favorite", h.MustBeAuthorizedMiddleware, h.RemoveFromFavorites)
			projectGroup.POST("/create-project", h.IsAdminMiddlware, h.CreateProject)
			projectGroup.DELETE("/:id", h.IsAdminMiddlware, h.DeleteProject)
			projectGroup.PUT("/:id", h.IsAdminMiddlware, h.UpdateProject)
		}

		movieGroup := authGroup.Group("/movies")
		{
			movieGroup.GET("/", h.GetFilteredMovies)
			movieGroup.GET("/:id", h.GetMovie)
		}

		seriesGroup := authGroup.Group("/series")
		{
			seriesGroup.GET("/", h.GetFilteredSeries)
			seriesGroup.GET("/:seriesID", h.GetSeries)

			seasonGroup := seriesGroup.Group("/seasons")
			{
				seasonGroup.GET("/:seasonNumber", h.GetSeason)

				episodeGroup := seriesGroup.Group("/episodes")
				{
					episodeGroup.GET("/:episodeID", h.GetEpisode)
				}
			}

		}

		genreGroup := authGroup.Group("/genres")
		{
			genreGroup.GET("/", h.GetAllGenres)
			genreGroup.GET("/:id", h.GetGenre)
			genreGroup.POST("/", h.IsAdminMiddlware, h.CreateGenre)
			genreGroup.DELETE("/:id", h.IsAdminMiddlware, h.DeleteGenre)
		}

		ageGroup := authGroup.Group("/ages")
		{
			ageGroup.GET("/", h.GetAllAgeCategories)
			ageGroup.GET("/:id", h.GetAgeCategory)
			ageGroup.POST("/", h.IsAdminMiddlware, h.CreateAgeCategory)
			ageGroup.DELETE("/:id", h.IsAdminMiddlware, h.DeleteAgeCategory)
		}
	}

	return router
}
