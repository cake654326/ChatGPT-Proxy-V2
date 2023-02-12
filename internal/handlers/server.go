package handlers

import (
	"github.com/acheong08/ChatGPT-V2/internal/api"
	"github.com/acheong08/ChatGPT-V2/internal/types"
	"github.com/gin-gonic/gin"
)

// Completions
func Completions(c *gin.Context) {
	// Map request body to CompletionRequest struct
	var request types.CompletionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	// Check if authorization header is present
	if c.GetHeader("Authorization") == "" {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	// Add authorization header to request
	request.Authorization = c.GetHeader("Authorization")
	// Check if prompt is present
	if request.Prompt == "" {
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// Send request to OpenAI API and stream data to client
	api.Send(request, c.Writer)
}
