package api

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"time"

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
	request.Stream = true // Temporary fix for OpenAI API
	if request.Paid {
		println("PAID")
		config.Model = "text-davinci-002-render-paid"
	} else {
		println("FREE")
		config.Model = "text-davinci-002-render"
	}
	body := map[string]interface{}{
		config.Mappings["model"]:            config.Model,
		config.Mappings["presence_penalty"]: request.PresencePenalty,
		config.Mappings["temperature"]:      request.Temperature,
		config.Mappings["top_p"]:            request.TopP,
		config.Mappings["stop"]:             request.Stop,
		config.Mappings["max_tokens"]:       request.MaxTokens,
		config.Mappings["stream"]:           request.Stream,
		config.Mappings["prompt"]:           request.Prompt,
	}
	// Create request
	req, err := http.NewRequest("POST", config.Endpoint, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
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
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	// Add body to request
	req.Body = io.NopCloser(bytes.NewReader(body_json))
	// Send request
	client := http.Client{
		Timeout: 360 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != 200 {
		c.JSON(503, gin.H{"error": "OpenAI error"})
		return
	}

	if request.Stream {
		// Use a buffer to store the response
		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if err != nil && err != io.EOF {
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			}
			if n == 0 {
				break
			}
			// Convert buf to string
			buf_str := string(buf[:n])
			// remove config.SecretModel from buf_str
			buf_str = regexp.MustCompile(config.SecretModel).ReplaceAllString(buf_str, "text-davinci-002-render")
			// Regex remove cmpl-6j6Ha2KTxZblH9BIu5FWhs1xUgpc3
			buf_str = regexp.MustCompile("cmpl-[a-zA-Z0-9]{29}").ReplaceAllString(buf_str, "...")
			// Regex replace "created": 1676206997 with "created": 0
			buf_str = regexp.MustCompile(`"created": [0-9]{10}`).ReplaceAllString(buf_str, `"created": 0`)
			// Make new buf from buf_str
			buf := []byte(buf_str)
			// Get new n from buf2
			n = len(buf)
			// Write the response chunk to the writer
			if _, err := writer.Write(buf[:n]); err != nil {
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			}
			// Flush the writer to ensure the response is sent immediately
			if f, ok := writer.(http.Flusher); ok {
				f.Flush()
			}
		}
	} else {
		// Read response body
		response_body := &bytes.Buffer{}
		_, err := response_body.ReadFrom(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}
		full_text := ""
		// Loop through each line of the response body choices finish_details is not null
		for {
			line, err := response_body.ReadString('\n')
			if err != nil && err != io.EOF {
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			}
			if line == "data: [DONE]" {
				break
			} else if line == "" {
				continue
			} else if line == "\n" {
				continue
			}
			// Remove the "data: " prefix
			line = line[6:]
			// Parse the line as JSON
			line_json := map[string]interface{}{}
			if json.Unmarshal([]byte(line), &line_json) != nil {
				c.JSON(500, gin.H{"error": "Internal server error"})
				return
			}
			// Look for line_json["choices"][0]["finish_details"]
			if line_json["choices"] != nil {
				if line_json["choices"].([]interface{})[0] != nil {
					// Check for text
					if line_json["choices"].([]interface{})[0].(map[string]interface{})["text"] != nil {
						// Append text to full_text
						full_text += line_json["choices"].([]interface{})[0].(map[string]interface{})["text"].(string)
					}
					if line_json["choices"].([]interface{})[0].(map[string]interface{})["finish_details"] != nil {
						response_body = bytes.NewBufferString(full_text)
						break
					}
				}
			}
		}
		c.Data(200, "application/json", response_body.Bytes())
	}
}
