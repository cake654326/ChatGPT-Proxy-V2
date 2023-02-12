package main

import (
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/acheong08/ChatGPT-V2/internal/handlers"
	"github.com/gin-gonic/gin"
)

var limit_middleware gin.HandlerFunc
var limit_store ratelimit.Store

func init() {
	limit_store = ratelimit.InMemoryStore(
		&ratelimit.InMemoryOptions{
			Rate:  time.Minute,
			Limit: 30,
		},
	)
	limit_middleware = ratelimit.RateLimiter(
		limit_store,
		&ratelimit.Options{
			ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
				c.JSON(
					429,
					gin.H{
						"message": "Too many requests",
					},
				)
				c.Abort()
			},
			KeyFunc: func(c *gin.Context) string {
				// Get Authorization header
				auth_header := c.GetHeader("Authorization")
				return auth_header
			},
		},
	)
}

func main() {
	handler := gin.Default()
	handler.Use(limit_middleware)
	handler.POST("/completions", handlers.Completions)
	handler.Run("127.0.0.1:10101")
}
