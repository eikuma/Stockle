package anthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	BaseURL = "https://api.anthropic.com/v1"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

type MessageRequest struct {
	Model     string    `json:"model"`
	MaxTokens int       `json:"max_tokens"`
	Messages  []Message `json:"messages"`
	System    string    `json:"system,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type MessageResponse struct {
	ID      string    `json:"id"`
	Type    string    `json:"type"`
	Role    string    `json:"role"`
	Content []Content `json:"content"`
	Model   string    `json:"model"`
	Usage   Usage     `json:"usage"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

func NewClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		apiKey:  apiKey,
		baseURL: BaseURL,
	}
}

func (c *Client) CreateMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/messages", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var msgResp MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &msgResp, nil
}