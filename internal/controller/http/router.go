package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "todo-clean/docs"
	"todo-clean/internal/usecase"
)

func NewRouter(
	logger zerolog.Logger,
	db *sql.DB,
	rdb *redis.Client,
	rl *RateLimiter,
	todoHandler *TodoHandler,
	authHandler *AuthHandler,
	authUC *usecase.AuthUseCase,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(requestID())
	r.Use(requestLogger(logger))
	r.Use(recovery(logger))
	r.Use(cors())

	r.GET("/healthz", healthCheck(db, rdb))
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")
	{
		v1.Use(rl.IPRateLimit())
		authHandler.Register(v1)

		protected := v1.Group("")
		protected.Use(jwtMiddleware(authUC))
		{
			todoHandler.Register(protected)
		}
	}

	return r
}

func healthCheck(db *sql.DB, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := gin.H{
			"status":   "ok",
			"postgres": "ok",
			"redis":    "ok",
		}
		httpStatus := http.StatusOK

		// check postgres
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			status["postgres"] = "unavailable"
			status["status"] = "degraded"
			httpStatus = http.StatusServiceUnavailable
		}

		// check redis with a throwaway client that bypasses the main pool
		if err := checkRedis(rdb.Options().Addr); err != nil {
			status["redis"] = "unavailable"
			status["status"] = "degraded"
			httpStatus = http.StatusServiceUnavailable
		}

		c.JSON(httpStatus, status)
	}
}

// checkRedis creates a fresh one-shot connection to verify Redis is reachable.
// A dedicated client bypasses the main pool which can mask connectivity issues
// by serving from cached connections.
func checkRedis(addr string) error {
	hc := redis.NewClient(&redis.Options{
		Addr:        addr,
		PoolSize:    1,
		DialTimeout: 2 * time.Second,
	})
	defer hc.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := hc.Ping(ctx).Err()
	fmt.Printf("checkRedis addr=%s err=%v\n", addr, err)
	return err
}
