package main

import (
	"os"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/acheong08/ChatGPT-V2/internal/api"
	"github.com/acheong08/ChatGPT-V2/internal/handlers"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	_ "github.com/go-redis/redis/v8"
)

var limit_middleware gin.HandlerFunc
var limit_store ratelimit.Store

func init() {
	limit_store = ratelimit.InMemoryStore(
		&ratelimit.InMemoryOptions{
			Rate:  time.Minute,
			Limit: 40,
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
				return c.ClientIP()
			},
		},
	)
}

func secret_auth(c *gin.Context) {
	if os.Getenv("SECRET") == "" {
		return
	}
	auth_header := c.GetHeader("Secret")
	if auth_header != os.Getenv("SECRET") {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		c.Abort()
		return
	}
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "10101"
	}
	handler := gin.Default()
	if !api.Config.Private {
		handler.Use(limit_middleware)
	}
	handler.Use(secret_auth)
	// Proxy all POST requests to https://apps.openai.com/api/
	handler.POST("/api/:path/", handlers.Proxy)
	// Proxy all GET requests to https://apps.openai.com/api/
	handler.GET("/api/:path/", handlers.Proxy)

	endless.ListenAndServe("127.0.0.1:"+PORT, handler)
}
