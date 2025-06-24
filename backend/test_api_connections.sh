#!/bin/bash

# AI API接続テストスクリプト
# 使用方法: ./test_api_connections.sh <groq_api_key> <claude_api_key>

set -e

echo "=== AI API接続テスト ==="
echo

# 引数チェック
if [ $# -ne 2 ]; then
    echo "使用方法: $0 <groq_api_key> <claude_api_key>"
    echo "例: $0 gsk_xxx sk-ant-xxx"
    exit 1
fi

GROQ_API_KEY=$1
CLAUDE_API_KEY=$2

# 1. Groq APIテスト
echo "1. Groq API接続テスト..."
echo "エンドポイント: https://api.groq.com/openai/v1/models"

GROQ_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  "https://api.groq.com/openai/v1/models" \
  -H "Authorization: Bearer $GROQ_API_KEY" \
  -H "Content-Type: application/json")

HTTP_CODE=$(echo "$GROQ_RESPONSE" | tail -n1 | cut -d: -f2)
RESPONSE_BODY=$(echo "$GROQ_RESPONSE" | sed '$d')

echo "HTTPステータス: $HTTP_CODE"
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Groq API接続成功"
    echo "利用可能モデル数: $(echo "$RESPONSE_BODY" | grep -o '"id"' | wc -l)"
    echo "サンプルモデル: $(echo "$RESPONSE_BODY" | grep -o '"id":"[^"]*"' | head -3)"
else
    echo "❌ Groq API接続失敗"
    echo "レスポンス: $RESPONSE_BODY"
fi
echo

# 2. Groq Chat Completions APIテスト
echo "2. Groq Chat Completions APIテスト..."
GROQ_CHAT_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  "https://api.groq.com/openai/v1/chat/completions" \
  -H "Authorization: Bearer $GROQ_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "llama-3.1-8b-instant",
    "messages": [{"role": "user", "content": "Hello, this is a test message. Please respond with JSON."}],
    "max_tokens": 100,
    "temperature": 0.1
  }')

HTTP_CODE=$(echo "$GROQ_CHAT_RESPONSE" | tail -n1 | cut -d: -f2)
RESPONSE_BODY=$(echo "$GROQ_CHAT_RESPONSE" | sed '$d')

echo "HTTPステータス: $HTTP_CODE"
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Groq Chat API接続成功"
    echo "レスポンス構造を確認:"
    echo "$RESPONSE_BODY" | python3 -m json.tool --indent 2 | head -20
else
    echo "❌ Groq Chat API接続失敗"
    echo "レスポンス: $RESPONSE_BODY"
fi
echo

# 3. Anthropic Claude APIテスト
echo "3. Anthropic Claude API接続テスト..."
echo "エンドポイント: https://api.anthropic.com/v1/messages"

CLAUDE_RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  "https://api.anthropic.com/v1/messages" \
  -H "x-api-key: $CLAUDE_API_KEY" \
  -H "Content-Type: application/json" \
  -H "anthropic-version: 2023-06-01" \
  -d '{
    "model": "claude-3-haiku-20240307",
    "max_tokens": 100,
    "messages": [{"role": "user", "content": "Hello, this is a test message. Please respond with JSON."}]
  }')

HTTP_CODE=$(echo "$CLAUDE_RESPONSE" | tail -n1 | cut -d: -f2)
RESPONSE_BODY=$(echo "$CLAUDE_RESPONSE" | sed '$d')

echo "HTTPステータス: $HTTP_CODE"
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Claude API接続成功"
    echo "レスポンス構造を確認:"
    echo "$RESPONSE_BODY" | python3 -m json.tool --indent 2 | head -20
else
    echo "❌ Claude API接続失敗"
    echo "レスポンス: $RESPONSE_BODY"
fi
echo

echo "=== テスト完了 ==="
echo "実際のAPIキーを使用して上記スクリプトを実行してください:"
echo "./test_api_connections.sh <your_groq_key> <your_claude_key>"