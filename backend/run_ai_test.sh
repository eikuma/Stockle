#!/bin/bash

# AI統合テスト実行スクリプト

echo "🚀 AI統合基盤テストを開始します"

# APIキーの確認
if [ -z "$GROQ_API_KEY" ] && [ -z "$ANTHROPIC_API_KEY" ]; then
    echo "❌ 環境変数 GROQ_API_KEY または ANTHROPIC_API_KEY を設定してください"
    echo ""
    echo "例:"
    echo "export GROQ_API_KEY=your_groq_api_key"
    echo "export ANTHROPIC_API_KEY=your_anthropic_api_key"
    echo ""
    echo "テストスキップ用の実行:"
    echo "export GROQ_API_KEY=dummy"
    echo "export ANTHROPIC_API_KEY=dummy"
    exit 1
fi

echo "🔧 Go依存関係を確認中..."
go mod tidy

echo "🏗️  プロジェクトをビルド中..."
go build ./...

if [ $? -ne 0 ]; then
    echo "❌ ビルドが失敗しました"
    exit 1
fi

echo "✅ ビルド成功"

echo "🧪 ユニットテストを実行中..."
go test ./internal/services -v

echo "🔗 AI統合テストを実行中..."
go run test_ai_integration.go

echo "🎉 AI統合基盤テスト完了"