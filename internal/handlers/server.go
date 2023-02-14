package handlers

import (
	"github.com/acheong08/ChatGPT-V2/internal/api"
	"github.com/gin-gonic/gin"
)

func Proxy(c *gin.Context) {
	// Send request to OpenAI API and stream data to client
	api.Proxy(c)
}
