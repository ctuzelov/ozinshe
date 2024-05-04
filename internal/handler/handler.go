package handler

import (
	"html/template"
	"log/slog"
	"ozinshe/internal/service"

	_ "ozinshe/docs"
	_ "ozinshe/internal/models"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service   *service.Service
	Log       *slog.Logger
	TempCache *template.Template
}

func New(service *service.Service, log *slog.Logger) *Handler {
	return &Handler{
		Service: service,
		Log:     log,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", h.HomePage)
	router.POST("/signup", h.SignUp)
	router.POST("/signin", h.SignIn)

	authGroup := router.Group("/", h.Middleware)
	{
		authGroup.GET("/signout", h.MustBeAuthorizedMiddleware, h.Signout)

		authGroup.GET("/users", h.IsAdminMiddlware, h.GetAllUsers)
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
			movieGroup.GET("/", h.GetAllMovies)
			movieGroup.GET("/:id", h.GetMovie)
		}

		seriesGroup := authGroup.Group("/series")
		{
			seriesGroup.GET("/", h.GetAllSeries)
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
