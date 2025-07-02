// test_ai_integration.go
// AI統合テスト用のスクリプト
// 実際のAPIキーを設定してテストを実行

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/eikuma/stockle/backend/internal/config"
	"github.com/eikuma/stockle/backend/internal/services"
)

func main() {
	fmt.Println("🚀 AI統合テスト開始")

	// 環境変数の確認
	if os.Getenv("GROQ_API_KEY") == "" {
		log.Println("⚠️  GROQ_API_KEY が設定されていません。Groq APIのテストをスキップします。")
	}

	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		log.Println("⚠️  ANTHROPIC_API_KEY が設定されていません。Claude APIのテストをスキップします。")
	}

	if os.Getenv("GROQ_API_KEY") == "" && os.Getenv("ANTHROPIC_API_KEY") == "" {
		log.Fatal("❌ 少なくとも1つのAPIキーを設定してください")
	}

	// 設定の初期化
	cfg := &config.AIConfig{
		GroqAPIKey:      os.Getenv("GROQ_API_KEY"),
		AnthropicAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		MaxRetries:      3,
		TimeoutSeconds:  30,
		RateLimit:       100,
	}

	// AIサービスの初期化
	aiService := services.NewAIService(cfg)

	// テスト用の記事内容
	testCases := []struct {
		name        string
		content     string
		title       string
		summaryType string
	}{
		{
			name:    "短文要約テスト",
			content: "人工知能（AI）技術の発展により、様々な分野での自動化が進んでいます。特に自然言語処理の分野では、大規模言語モデルが注目を集めており、文章の要約、翻訳、質問応答などのタスクで高い性能を示しています。",
			title:   "AI技術の発展",
			summaryType: "short",
		},
		{
			name:    "中文要約テスト",
			content: "近年、リモートワークが急速に普及し、働き方に大きな変化をもたらしています。企業は柔軟な働き方を支援するため、クラウドベースのツールやコミュニケーションプラットフォームの導入を進めています。この変化により、ワークライフバランスの改善や地理的制約の解消などのメリットが生まれる一方で、チームコミュニケーションの課題や労働時間管理の複雑化などの新たな問題も浮き彫りになっています。",
			title:   "リモートワークの普及と課題",
			summaryType: "medium",
		},
	}

	ctx := context.Background()

	// 各テストケースを実行
	for _, tc := range testCases {
		fmt.Printf("\n📝 %s を実行中...\n", tc.name)

		req := &services.SummaryRequest{
			Content:     tc.content,
			Title:       tc.title,
			URL:         "https://example.com/test-article",
			Language:    "ja",
			SummaryType: tc.summaryType,
		}

		result, err := aiService.GenerateSummary(ctx, req)
		if err != nil {
			fmt.Printf("❌ %s 失敗: %v\n", tc.name, err)
			continue
		}

		// 結果の表示
		fmt.Printf("✅ %s 成功\n", tc.name)
		fmt.Printf("   プロバイダー: %s\n", result.Provider)
		fmt.Printf("   要約: %s\n", result.Summary)
		fmt.Printf("   信頼度: %.2f\n", result.Confidence)
		fmt.Printf("   文字数: %d\n", result.WordCount)
		fmt.Printf("   生成時刻: %s\n", result.GeneratedAt.Format("2006-01-02 15:04:05"))
	}

	fmt.Println("\n🎉 AI統合テスト完了")
}