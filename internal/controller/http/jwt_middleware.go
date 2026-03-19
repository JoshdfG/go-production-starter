package http

import (
	"net/http"
	"strings"

	"todo-clean/internal/usecase"

	"github.com/gin-gonic/gin"
)

func jwtMiddleware(authUC *usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error: "authorization header required",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error: "authorization header format: Bearer <token>",
			})
			return
		}

		userID, err := authUC.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
				Error: "invalid or expired token",
			})
			return
		}

		// store userID in context for handlers to use
		c.Set("userID", userID)
		c.Next()
	}
}
