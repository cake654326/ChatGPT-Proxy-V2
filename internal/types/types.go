package types

type Config struct {
	Endpoint string `json:"endpoint"`
	Model    string `json:"model"`
}

type CompletionRequest struct {
	Prompt          string   `json:"prompt"`
	MaxTokens       int      `json:"max_tokens"`
	Temperature     float32  `json:"temperature"`
	TopP            float32  `json:"top_p"`
	Stop            []string `json:"stop"`
	PresencePenalty float32  `json:"presence_penalty"`
	Authorization   string   `json:"authorization"`
}
