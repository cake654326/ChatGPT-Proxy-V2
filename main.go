package main

import (
	"os"

	"github.com/acheong08/ChatGPT-V2/internal/handlers"
	"github.com/gin-gonic/gin"
)

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
	handler := gin.Default()
	handler.Use(secret_auth)
	handler.POST("/completions", handlers.Completions)
	handler.Run("127.0.0.1:10101")
}
