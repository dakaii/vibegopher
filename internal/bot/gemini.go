package bot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type GeminiClient struct {
	apiKey     string
	httpClient *http.Client
	model      string
}

func NewGeminiClient(apiKey string) *GeminiClient {
	return &GeminiClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 45 * time.Second,
		},
		model: "gemini-2.0-flash",
	}
}

type geminiRequest struct {
	SystemInstruction *geminiContent `json:"system_instruction,omitempty"`
	Contents          []geminiContent `json:"contents"`
	GenerationConfig  geminiGenConfig `json:"generationConfig"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiGenConfig struct {
	MaxOutputTokens int     `json:"maxOutputTokens"`
	Temperature     float64 `json:"temperature"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *GeminiClient) GenerateComment(ctx context.Context, userPrompt string) (string, error) {
	if c.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not configured")
	}

	body := geminiRequest{
		SystemInstruction: &geminiContent{Parts: []geminiPart{{Text: PersonaSystemPrompt}}},
		Contents: []geminiContent{
			{Parts: []geminiPart{{Text: userPrompt}}},
		},
		GenerationConfig: geminiGenConfig{
			MaxOutputTokens: 200,
			Temperature:     0.8,
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		c.model,
		c.apiKey,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("gemini http %d: %s", resp.StatusCode, truncate(string(respBody), 300))
	}

	var parsed geminiResponse
	if err := json.Unmarshal(respBody, &parsed); err != nil {
		return "", err
	}
	if parsed.Error != nil {
		return "", fmt.Errorf("gemini: %s", parsed.Error.Message)
	}
	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty response")
	}

	text := strings.TrimSpace(parsed.Candidates[0].Content.Parts[0].Text)
	text = strings.Trim(text, `"'`)
	if len(text) > 280 {
		text = text[:277] + "..."
	}
	if text == "" {
		return "", fmt.Errorf("gemini returned blank comment")
	}
	return text, nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
