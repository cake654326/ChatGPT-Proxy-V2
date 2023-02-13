package types

type Config struct {
	Endpoint    string            `json:"endpoint"`
	Model       string            `json:"model"`
	Mappings    map[string]string `json:"mappings"`
	SecretModel string            `json:"secret_model"`
	Private     bool              `json:"private"`
}

type CompletionRequest struct {
	Prompt          string   `json:"prompt"`
	MaxTokens       int      `json:"max_tokens"`
	Temperature     float32  `json:"temperature"`
	TopP            float32  `json:"top_p"`
	Stop            []string `json:"stop"`
	PresencePenalty float32  `json:"presence_penalty"`
	Authorization   string   `json:"authorization"`
	Stream          bool     `json:"stream"`
	Paid            bool     `json:"paid"`
}
