package http

import (
	_ "todo-clean/docs"
	"todo-clean/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	logger zerolog.Logger,
	todoHandler *TodoHandler,
	authHandler *AuthHandler,
	authUC *usecase.AuthUseCase,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(requestLogger(logger))
	r.Use(recovery(logger))
	r.Use(cors())

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")
	{
		// public routes — no auth required
		authHandler.Register(v1)

		// protected routes — JWT required
		protected := v1.Group("")
		protected.Use(jwtMiddleware(authUC))
		{
			todoHandler.Register(protected)
		}
	}

	return r
}
