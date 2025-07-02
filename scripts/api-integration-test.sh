#!/bin/bash

# Stockle API統合テストスクリプト
# 実行前に必要: jq, curl
# 使用方法: ./api-integration-test.sh [API_URL]

set -e

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# デフォルト設定
API_URL="${1:-http://localhost:8080}"
BASE_URL="${API_URL}/api/v1"

# テスト用データ
TEST_EMAIL="test_$(date +%s)@example.com"
TEST_PASSWORD="testpass123"
TEST_DISPLAY_NAME="Test User"
TEST_ARTICLE_URL="https://example.com/article-$(date +%s)"

# グローバル変数
ACCESS_TOKEN=""
REFRESH_TOKEN=""
USER_ID=""
ARTICLE_ID=""
CATEGORY_ID=""

# 関数定義
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}→ $1${NC}"
}

check_response() {
    local response=$1
    local expected_status=$2
    local actual_status=$(echo "$response" | tail -n1)
    
    if [ "$actual_status" = "$expected_status" ]; then
        return 0
    else
        return 1
    fi
}

# テスト開始
echo "================================================"
echo "    Stockle API統合テスト"
echo "================================================"
echo "API URL: $BASE_URL"
echo "------------------------------------------------"

# 1. ヘルスチェック
print_info "1. ヘルスチェック"
HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" "$BASE_URL/health")
if check_response "$HEALTH_RESPONSE" "200"; then
    print_success "APIサーバーは正常に動作しています"
else
    print_error "APIサーバーに接続できません"
    exit 1
fi

# 2. ユーザー登録
print_info "2. ユーザー登録"
REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/register" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$TEST_EMAIL\",
        \"password\": \"$TEST_PASSWORD\",
        \"display_name\": \"$TEST_DISPLAY_NAME\"
    }")

if check_response "$REGISTER_RESPONSE" "201"; then
    USER_DATA=$(echo "$REGISTER_RESPONSE" | head -n-1)
    USER_ID=$(echo "$USER_DATA" | jq -r '.user.id')
    print_success "ユーザー登録成功 (ID: $USER_ID)"
else
    print_error "ユーザー登録失敗"
    echo "$REGISTER_RESPONSE"
fi

# 3. ログイン
print_info "3. ログイン"
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{
        \"email\": \"$TEST_EMAIL\",
        \"password\": \"$TEST_PASSWORD\"
    }")

if check_response "$LOGIN_RESPONSE" "200"; then
    LOGIN_DATA=$(echo "$LOGIN_RESPONSE" | head -n-1)
    ACCESS_TOKEN=$(echo "$LOGIN_DATA" | jq -r '.tokens.access_token')
    REFRESH_TOKEN=$(echo "$LOGIN_DATA" | jq -r '.tokens.refresh_token')
    print_success "ログイン成功"
else
    print_error "ログイン失敗"
    echo "$LOGIN_RESPONSE"
    exit 1
fi

# 4. プロフィール取得
print_info "4. プロフィール取得"
PROFILE_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/users/me" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

if check_response "$PROFILE_RESPONSE" "200"; then
    print_success "プロフィール取得成功"
else
    print_error "プロフィール取得失敗"
fi

# 5. カテゴリ作成
print_info "5. カテゴリ作成"
CATEGORY_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/categories" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "{
        \"name\": \"Technology\",
        \"color\": \"#FF5722\"
    }")

if check_response "$CATEGORY_RESPONSE" "201"; then
    CATEGORY_DATA=$(echo "$CATEGORY_RESPONSE" | head -n-1)
    CATEGORY_ID=$(echo "$CATEGORY_DATA" | jq -r '.id')
    print_success "カテゴリ作成成功 (ID: $CATEGORY_ID)"
else
    print_error "カテゴリ作成失敗"
fi

# 6. 記事保存
print_info "6. 記事保存"
ARTICLE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/articles" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $ACCESS_TOKEN" \
    -d "{
        \"url\": \"$TEST_ARTICLE_URL\",
        \"category_id\": \"$CATEGORY_ID\",
        \"tags\": [\"test\", \"integration\"]
    }")

if check_response "$ARTICLE_RESPONSE" "201"; then
    ARTICLE_DATA=$(echo "$ARTICLE_RESPONSE" | head -n-1)
    ARTICLE_ID=$(echo "$ARTICLE_DATA" | jq -r '.id')
    print_success "記事保存成功 (ID: $ARTICLE_ID)"
else
    print_error "記事保存失敗"
    echo "$ARTICLE_RESPONSE"
fi

# 7. 記事一覧取得
print_info "7. 記事一覧取得"
ARTICLES_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/articles?limit=10" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

if check_response "$ARTICLES_RESPONSE" "200"; then
    ARTICLES_DATA=$(echo "$ARTICLES_RESPONSE" | head -n-1)
    ARTICLE_COUNT=$(echo "$ARTICLES_DATA" | jq '.total')
    print_success "記事一覧取得成功 (記事数: $ARTICLE_COUNT)"
else
    print_error "記事一覧取得失敗"
fi

# 8. 記事詳細取得
if [ -n "$ARTICLE_ID" ]; then
    print_info "8. 記事詳細取得"
    ARTICLE_DETAIL_RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$BASE_URL/articles/$ARTICLE_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if check_response "$ARTICLE_DETAIL_RESPONSE" "200"; then
        print_success "記事詳細取得成功"
    else
        print_error "記事詳細取得失敗"
    fi
fi

# 9. 要約生成リクエスト
if [ -n "$ARTICLE_ID" ]; then
    print_info "9. 要約生成リクエスト"
    SUMMARY_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/articles/$ARTICLE_ID/summary" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if check_response "$SUMMARY_RESPONSE" "202"; then
        print_success "要約生成リクエスト成功"
        # 実際の環境では要約生成完了を待つ必要があります
    else
        print_error "要約生成リクエスト失敗"
    fi
fi

# 10. 記事ステータス更新
if [ -n "$ARTICLE_ID" ]; then
    print_info "10. 記事ステータス更新"
    UPDATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/articles/$ARTICLE_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $ACCESS_TOKEN" \
        -d "{
            \"status\": \"read\"
        }")
    
    if check_response "$UPDATE_RESPONSE" "200"; then
        print_success "記事ステータス更新成功"
    else
        print_error "記事ステータス更新失敗"
    fi
fi

# 11. トークンリフレッシュ
print_info "11. トークンリフレッシュ"
REFRESH_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/refresh" \
    -H "Content-Type: application/json" \
    -d "{
        \"refresh_token\": \"$REFRESH_TOKEN\"
    }")

if check_response "$REFRESH_RESPONSE" "200"; then
    REFRESH_DATA=$(echo "$REFRESH_RESPONSE" | head -n-1)
    NEW_ACCESS_TOKEN=$(echo "$REFRESH_DATA" | jq -r '.tokens.access_token')
    print_success "トークンリフレッシュ成功"
else
    print_error "トークンリフレッシュ失敗"
fi

# 12. 記事削除
if [ -n "$ARTICLE_ID" ]; then
    print_info "12. 記事削除"
    DELETE_RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "$BASE_URL/articles/$ARTICLE_ID" \
        -H "Authorization: Bearer $ACCESS_TOKEN")
    
    if check_response "$DELETE_RESPONSE" "204"; then
        print_success "記事削除成功"
    else
        print_error "記事削除失敗"
    fi
fi

# 13. ログアウト
print_info "13. ログアウト"
LOGOUT_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/auth/logout" \
    -H "Authorization: Bearer $ACCESS_TOKEN")

if check_response "$LOGOUT_RESPONSE" "200"; then
    print_success "ログアウト成功"
else
    print_error "ログアウト失敗"
fi

echo "------------------------------------------------"
echo "API統合テスト完了"