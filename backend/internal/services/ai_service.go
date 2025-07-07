package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/private/Stockle/backend/internal/config"
	"github.com/private/Stockle/backend/pkg/anthropic"
	"github.com/private/Stockle/backend/pkg/groq"
)

type AIService struct {
	groqClient      *groq.Client
	anthropicClient *anthropic.Client
	config          *config.AIConfig
}

type SummaryRequest struct {
	Content     string
	Title       string
	URL         string
	Language    string
	SummaryType string // "short", "medium", "long"
}

type SummaryResponse struct {
	Summary      string
	Confidence   float64
	Provider     string
	GeneratedAt  time.Time
	ModelVersion string
	WordCount    int
}

func NewAIService(cfg *config.AIConfig) *AIService {
	return &AIService{
		groqClient:      groq.NewClient(cfg.GroqAPIKey),
		anthropicClient: anthropic.NewClient(cfg.AnthropicAPIKey),
		config:          cfg,
	}
}

func (s *AIService) GenerateSummary(ctx context.Context, req *SummaryRequest) (*SummaryResponse, error) {
	// Groqを第一選択として使用
	summary, err := s.generateWithGroq(ctx, req)
	if err != nil {
		// Groqが失敗した場合はClaude APIにフォールバック
		summary, err = s.generateWithClaude(ctx, req)
		if err != nil {
			return nil, fmt.Errorf("all AI providers failed: %w", err)
		}
	}

	return summary, nil
}

func (s *AIService) generateWithGroq(ctx context.Context, req *SummaryRequest) (*SummaryResponse, error) {
	prompt := s.buildPrompt(req)

	groqReq := &groq.ChatCompletionRequest{
		Model: "llama3-8b-8192",
		Messages: []groq.Message{
			{
				Role:    "system",
				Content: s.getSystemPrompt(req.SummaryType),
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens:   500,
		Temperature: 0.3,
	}

	resp, err := s.groqClient.CreateChatCompletion(ctx, groqReq)
	if err != nil {
		return nil, fmt.Errorf("groq API error: %w", err)
	}

	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	confidence := s.calculateConfidence(summary, req.Content)

	return &SummaryResponse{
		Summary:      summary,
		Confidence:   confidence,
		Provider:     "groq",
		GeneratedAt:  time.Now(),
		ModelVersion: resp.Model,
		WordCount:    len(strings.Fields(summary)),
	}, nil
}

func (s *AIService) generateWithClaude(ctx context.Context, req *SummaryRequest) (*SummaryResponse, error) {
	prompt := s.buildPrompt(req)

	claudeReq := &anthropic.MessageRequest{
		Model:     "claude-3-haiku-20240307",
		MaxTokens: 500,
		Messages: []anthropic.Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		System: s.getSystemPrompt(req.SummaryType),
	}

	resp, err := s.anthropicClient.CreateMessage(ctx, claudeReq)
	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	summary := strings.TrimSpace(resp.Content[0].Text)
	confidence := s.calculateConfidence(summary, req.Content)

	return &SummaryResponse{
		Summary:      summary,
		Confidence:   confidence,
		Provider:     "claude",
		GeneratedAt:  time.Now(),
		ModelVersion: resp.Model,
		WordCount:    len(strings.Fields(summary)),
	}, nil
}

func (s *AIService) buildPrompt(req *SummaryRequest) string {
	return fmt.Sprintf(`記事タイトル: %s

記事URL: %s

記事内容:
%s

上記の記事を%sで要約してください。`, req.Title, req.URL, req.Content, s.getSummaryTypeDescription(req.SummaryType))
}

func (s *AIService) getSystemPrompt(summaryType string) string {
	switch summaryType {
	case "short":
		return "あなたは記事の要約を作成する専門家です。記事の主要なポイントを50-100文字で簡潔に要約してください。読者が記事の概要を素早く理解できるようにしてください。"
	case "long":
		return "あなたは記事の要約を作成する専門家です。記事の詳細な内容を500-800文字で包括的に要約してください。重要なポイント、詳細、結論を含めて、読者が記事を読まなくても内容を十分理解できるようにしてください。"
	default: // medium
		return "あなたは記事の要約を作成する専門家です。記事の主要なポイントを200-300文字で要約してください。重要な情報を含みつつ、読みやすい長さにまとめてください。"
	}
}

func (s *AIService) getSummaryTypeDescription(summaryType string) string {
	switch summaryType {
	case "short":
		return "50-100文字の短い要約"
	case "long":
		return "500-800文字の詳細な要約"
	default:
		return "200-300文字の標準的な要約"
	}
}

func (s *AIService) calculateConfidence(summary, originalContent string) float64 {
	// 簡単な信頼度計算（実際にはより複雑なロジックが必要）
	if len(summary) < 50 {
		return 0.6
	}
	if len(summary) > len(originalContent) {
		return 0.4
	}
	return 0.8
}