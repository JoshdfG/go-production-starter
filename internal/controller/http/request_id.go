package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// requestID injects a unique ID into every request.
// If the client sends X-Request-ID we honour it, otherwise we generate one.
// The ID is attached to the Gin context and returned in the response header
// so clients can correlate their request with your logs.
func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Set(RequestIDHeader, id)
		c.Header(RequestIDHeader, id)
		c.Next()
	}
}
