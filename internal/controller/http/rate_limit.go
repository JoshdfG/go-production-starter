package http

import (
	"fmt"
	"net/http"
	"time"

	"todo-clean/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RateLimiter holds limiters for different contexts.
type RateLimiter struct {
	ipLimiter   *limiter.Limiter
	userLimiter *limiter.Limiter
}

// NewRateLimiter creates two limiters:
// - IP limiter:   100 requests per minute for unauthenticated traffic
// - User limiter: 1000 requests per minute for authenticated users.
func NewRateLimiter(rdb *redis.Client, cfg config.RateLimitConfig) (*RateLimiter, error) {
	store, err := redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{
		Prefix:   "rate",
		MaxRetry: 3,
	})
	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		ipLimiter: limiter.New(store, limiter.Rate{
			Period: time.Minute,
			Limit:  int64(cfg.IPLimit),
		}),
		userLimiter: limiter.New(store, limiter.Rate{
			Period: time.Minute,
			Limit:  int64(cfg.UserLimit),
		}),
	}, nil
}

// IPRateLimit limits by client IP — used on public routes.
func (rl *RateLimiter) IPRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, err := rl.ipLimiter.Get(c.Request.Context(), c.ClientIP())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
				Error: "rate limiter error",
			})
			return
		}

		// set standard rate limit headers
		c.Header("X-RateLimit-Limit", formatInt(ctx.Limit))
		c.Header("X-RateLimit-Remaining", formatInt(ctx.Remaining))
		c.Header("X-RateLimit-Reset", formatInt(ctx.Reset))

		if ctx.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{
				Error: "rate limit exceeded — try again later",
			})
			return
		}
		c.Next()
	}
}

// UserRateLimit limits by authenticated user ID — used on protected routes.
// Falls back to IP if no userID is in context.
func (rl *RateLimiter) UserRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()
		if userID, exists := c.Get("userID"); exists {
			key = "user:" + userID.(string)
		}

		ctx, err := rl.userLimiter.Get(c.Request.Context(), key)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
				Error: "rate limiter error",
			})
			return
		}

		c.Header("X-RateLimit-Limit", formatInt(ctx.Limit))
		c.Header("X-RateLimit-Remaining", formatInt(ctx.Remaining))
		c.Header("X-RateLimit-Reset", formatInt(ctx.Reset))

		if ctx.Reached {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, ErrorResponse{
				Error: "rate limit exceeded — try again later",
			})
			return
		}
		c.Next()
	}
}

func formatInt(i int64) string {
	return fmt.Sprintf("%d", i)
}
