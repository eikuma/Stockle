package services

import (
	"context"
	"os"
	"testing"

	"github.com/private/Stockle/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAIService_GenerateSummary(t *testing.T) {
	// テスト環境でAPIキーが設定されていない場合はスキップ
	if os.Getenv("GROQ_API_KEY") == "" || os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("API keys not set, skipping integration test")
	}

	cfg, err := config.NewAIConfig()
	require.NoError(t, err)

	aiService := NewAIService(cfg)

	tests := []struct {
		name        string
		request     *SummaryRequest
		expectError bool
	}{
		{
			name: "正常な要約生成（短文）",
			request: &SummaryRequest{
				Content: `
				人工知能（AI）技術の発展により、様々な分野での自動化が進んでいます。
				特に自然言語処理の分野では、大規模言語モデルが注目を集めており、
				文章の要約、翻訳、質問応答などのタスクで高い性能を示しています。
				これらの技術は、情報処理の効率化や新しいサービスの創出に貢献しています。
				`,
				Title:       "AI技術の発展について",
				URL:         "https://example.com/ai-development",
				Language:    "ja",
				SummaryType: "short",
			},
			expectError: false,
		},
		{
			name: "正常な要約生成（中文）",
			request: &SummaryRequest{
				Content: `
				近年、リモートワークが急速に普及し、働き方に大きな変化をもたらしています。
				企業は柔軟な働き方を支援するため、クラウドベースのツールやコミュニケーション
				プラットフォームの導入を進めています。この変化により、ワークライフバランスの
				改善や地理的制約の解消などのメリットが生まれる一方で、チームコミュニケーションの
				課題や労働時間管理の複雑化などの新たな問題も浮き彫りになっています。
				`,
				Title:       "リモートワークの普及と課題",
				URL:         "https://example.com/remote-work",
				Language:    "ja",
				SummaryType: "medium",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			
			result, err := aiService.GenerateSummary(ctx, tt.request)
			
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.Summary)
				assert.Greater(t, result.Confidence, 0.0)
				assert.NotEmpty(t, result.Provider)
				assert.Greater(t, result.WordCount, 0)
				
				t.Logf("Provider: %s", result.Provider)
				t.Logf("Summary: %s", result.Summary)
				t.Logf("Confidence: %.2f", result.Confidence)
				t.Logf("Word Count: %d", result.WordCount)
			}
		})
	}
}

func TestAIService_ProviderFallback(t *testing.T) {
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("ANTHROPIC_API_KEY not set, skipping fallback test")
	}

	// Groq APIキーを無効にしてフォールバック機能をテスト
	cfg := &config.AIConfig{
		GroqAPIKey:      "invalid-key",
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		MaxRetries:      3,
		TimeoutSeconds:  30,
		RateLimit:       100,
	}

	aiService := NewAIService(cfg)

	request := &SummaryRequest{
		Content:     "テスト用の短い記事内容です。AIの要約機能をテストしています。",
		Title:       "テスト記事",
		URL:         "https://example.com/test",
		Language:    "ja",
		SummaryType: "short",
	}

	ctx := context.Background()
	result, err := aiService.GenerateSummary(ctx, request)

	// Groqが失敗してClaudeにフォールバックした場合
	if err == nil {
		assert.NotNil(t, result)
		assert.Equal(t, "claude", result.Provider)
		t.Logf("Fallback successful - Provider: %s", result.Provider)
		t.Logf("Summary: %s", result.Summary)
	} else {
		// 両方のプロバイダーが失敗した場合
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "all AI providers failed")
		t.Logf("Both providers failed as expected: %v", err)
	}
}