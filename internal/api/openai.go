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
func Send(request types.CompletionRequest, writer gin.ResponseWriter) {
	// Create HTTP headers
	headers := http.Header{
		"Authorization": []string{request.Authorization},
		"Content-Type":  []string{"application/json"},
	}
	// Create body JSON
	body := map[string]interface{}{
		"model":            config.Model,
		"presence_penalty": request.PresencePenalty,
		"temperature":      request.Temperature,
		"top_p":            request.TopP,
		"stop":             request.Stop,
		"max_tokens":       request.MaxTokens,
		"stream":           true,
		"query":            request.Prompt,
	}
	// Create request
	req, err := http.NewRequest("POST", config.Endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Add headers to request
	req.Header = headers
	// Build request body
	body_json, err := json.Marshal(body)
	if err != nil {
		log.Fatal(err)
	}
	// Add body to request
	req.Body = io.NopCloser(bytes.NewReader(body_json))
	// Send request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Use a buffer to store the response
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if n == 0 {
			break
		}
		// Write the response chunk to the writer
		if _, err := writer.Write(buf[:n]); err != nil {
			log.Fatal(err)
		}
		// Flush the writer to ensure the response is sent immediately
		if f, ok := writer.(http.Flusher); ok {
			f.Flush()
		}
	}
}
