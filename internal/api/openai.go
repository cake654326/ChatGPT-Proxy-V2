package api

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/acheong08/ChatGPT-V2/internal/types"
	"github.com/gin-gonic/gin"
)

var (
	//go:embed config.json
	config_file []byte
	config      types.Config
)

// config returns the config.json file as a Config struct.
func init() {
	config = types.Config{}
	if json.Unmarshal(config_file, &config) != nil {
		log.Fatal("Error unmarshalling config.json")
	}
}

// Returns a stream of completions from the OpenAI API.
func Send(request types.CompletionRequest, writer gin.ResponseWriter, c *gin.Context) {
	// Create HTTP headers
	headers := http.Header{
		"Authorization": []string{request.Authorization},
		"Content-Type":  []string{"application/json"},
	}
	// Create body JSON
	body := map[string]interface{}{
		config.Mappings["model"]:            config.Model,
		config.Mappings["presence_penalty"]: request.PresencePenalty,
		config.Mappings["temperature"]:      request.Temperature,
		config.Mappings["top_p"]:            request.TopP,
		config.Mappings["stop"]:             request.Stop,
		config.Mappings["max_tokens"]:       request.MaxTokens,
		config.Mappings["stream"]:           true,
		config.Mappings["prompt"]:           request.Prompt,
	}
	// Create request
	req, err := http.NewRequest("POST", config.Endpoint, nil)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return

	}
	// Set timeout as 360 seconds
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Keep-Alive", "360")
	// Add headers to request
	req.Header = headers
	// Build request body
	body_json, err := json.Marshal(body)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	// Add body to request
	req.Body = io.NopCloser(bytes.NewReader(body_json))
	// Send request
	client := http.Client{
		Timeout: 360,
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		c.JSON(resp.StatusCode, gin.H{"message": "Invalid request"})
		return
	}

	// Use a buffer to store the response
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			c.JSON(500, gin.H{"message": "Internal server error"})
			return
		}
		if n == 0 {
			break
		}
		// Write the response chunk to the writer
		if _, err := writer.Write(buf[:n]); err != nil {
			c.JSON(500, gin.H{"message": "Internal server error"})
			return
		}
		// Flush the writer to ensure the response is sent immediately
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}
