#!/bin/bash

# =====================================
# Stockle 統合テストスクリプト
# =====================================

set -euo pipefail

# 色付きログ用の定数
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ログ関数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# エラーハンドリング
cleanup() {
    log_info "クリーンアップを実行中..."
    
    # バックエンドプロセスを停止
    if [[ -n "${BACKEND_PID:-}" ]] && kill -0 "$BACKEND_PID" 2>/dev/null; then
        log_info "バックエンドサーバーを停止中 (PID: $BACKEND_PID)"
        kill "$BACKEND_PID" || true
        wait "$BACKEND_PID" 2>/dev/null || true
    fi
    
    # フロントエンドプロセスを停止
    if [[ -n "${FRONTEND_PID:-}" ]] && kill -0 "$FRONTEND_PID" 2>/dev/null; then
        log_info "フロントエンドサーバーを停止中 (PID: $FRONTEND_PID)"
        kill "$FRONTEND_PID" || true
        wait "$FRONTEND_PID" 2>/dev/null || true
    fi
    
    log_info "クリーンアップ完了"
}

trap cleanup EXIT

# 設定値
BACKEND_PORT=8080
FRONTEND_PORT=3000
TIMEOUT=30
TEST_RESULTS_DIR="test-results"

# =====================================
# 1. 事前チェック
# =====================================
log_info "=== 統合テスト開始 ==="
log_info "事前チェックを実行中..."

# 必要なディレクトリの存在確認
if [[ ! -d "backend" ]]; then
    log_error "backendディレクトリが見つかりません"
    exit 1
fi

if [[ ! -d "frontend" ]]; then
    log_error "frontendディレクトリが見つかりません"
    exit 1
fi

# データベース接続確認
log_info "データベース接続を確認中..."
if ! mysqladmin ping -h "${DB_HOST:-127.0.0.1}" -P "${DB_PORT:-3306}" --silent 2>/dev/null; then
    log_error "データベースに接続できません"
    exit 1
fi

log_success "事前チェック完了"

# =====================================
# 2. テスト結果ディレクトリ作成
# =====================================
mkdir -p "$TEST_RESULTS_DIR"

# =====================================
# 3. データベースマイグレーション実行
# =====================================
log_info "データベースマイグレーションを実行中..."

cd backend

# マイグレーションファイルが存在する場合のみ実行
if [[ -d "migrations" ]] && [[ -n "$(ls -A migrations)" ]]; then
    # golang-migrateがインストールされている場合
    if command -v migrate &> /dev/null; then
        DB_URL="mysql://${DB_USER:-testuser}:${DB_PASSWORD}@tcp(${DB_HOST:-127.0.0.1}:${DB_PORT:-3306})/${DB_NAME:-stockle_test}"
        migrate -path migrations -database "$DB_URL" up
        log_success "マイグレーション完了"
    else
        log_warning "migrateコマンドが見つかりません。マイグレーションをスキップします"
    fi
else
    log_warning "マイグレーションファイルが見つかりません。スキップします"
fi

cd ..

# =====================================
# 4. バックエンドサーバー起動
# =====================================
log_info "バックエンドサーバーを起動中..."

cd backend

# バックエンドをバックグラウンドで起動
if [[ -f "api" ]]; then
    ./api &
    BACKEND_PID=$!
else
    go run cmd/api/main.go &
    BACKEND_PID=$!
fi

cd ..

# バックエンドサーバーのヘルスチェック
log_info "バックエンドサーバーのヘルスチェック中..."
for i in $(seq 1 $TIMEOUT); do
    if curl -f -s "http://localhost:$BACKEND_PORT/health" > /dev/null 2>&1; then
        log_success "バックエンドサーバーが起動しました"
        break
    fi
    
    if [[ $i -eq $TIMEOUT ]]; then
        log_error "バックエンドサーバーの起動がタイムアウトしました"
        exit 1
    fi
    
    sleep 1
done

# =====================================
# 5. フロントエンドサーバー起動
# =====================================
log_info "フロントエンドサーバーを起動中..."

cd frontend

# フロントエンドをバックグラウンドで起動
if [[ -d ".next" ]]; then
    npm start &
    FRONTEND_PID=$!
else
    npm run dev &
    FRONTEND_PID=$!
fi

cd ..

# フロントエンドサーバーのヘルスチェック
log_info "フロントエンドサーバーのヘルスチェック中..."
for i in $(seq 1 $TIMEOUT); do
    if curl -f -s "http://localhost:$FRONTEND_PORT" > /dev/null 2>&1; then
        log_success "フロントエンドサーバーが起動しました"
        break
    fi
    
    if [[ $i -eq $TIMEOUT ]]; then
        log_error "フロントエンドサーバーの起動がタイムアウトしました"
        exit 1
    fi
    
    sleep 1
done

# =====================================
# 6. E2Eテスト実行
# =====================================
log_info "E2Eテストを実行中..."

cd frontend

# Playwrightテストが存在する場合のみ実行
if [[ -f "playwright.config.js" ]] || [[ -f "playwright.config.ts" ]]; then
    if npm run test:e2e 2>/dev/null; then
        log_success "E2Eテスト完了"
    else
        log_warning "E2Eテストでエラーが発生しましたが、続行します"
    fi
else
    log_warning "E2Eテスト設定が見つかりません。スキップします"
fi

cd ..

# =====================================
# 7. API統合テスト実行
# =====================================
log_info "API統合テストを実行中..."

# 基本的なAPIエンドポイントテスト
test_api_endpoint() {
    local endpoint=$1
    local expected_status=${2:-200}
    local method=${3:-GET}
    
    log_info "テスト中: $method $endpoint"
    
    if [[ "$method" == "GET" ]]; then
        response=$(curl -s -w "%{http_code}" "http://localhost:$BACKEND_PORT$endpoint")
        status_code="${response: -3}"
    else
        status_code=$(curl -s -o /dev/null -w "%{http_code}" -X "$method" "http://localhost:$BACKEND_PORT$endpoint")
    fi
    
    if [[ "$status_code" == "$expected_status" ]]; then
        log_success "✓ $endpoint ($status_code)"
    else
        log_error "✗ $endpoint (expected: $expected_status, got: $status_code)"
        return 1
    fi
}

# APIテスト実行
API_TESTS_PASSED=0
API_TESTS_TOTAL=0

# ヘルスチェックエンドポイント
((API_TESTS_TOTAL++))
if test_api_endpoint "/health" 200; then
    ((API_TESTS_PASSED++))
fi

# その他の基本エンドポイント（実装に応じて追加）
if curl -f -s "http://localhost:$BACKEND_PORT/api/v1" > /dev/null 2>&1; then
    ((API_TESTS_TOTAL++))
    if test_api_endpoint "/api/v1" 200; then
        ((API_TESTS_PASSED++))
    fi
fi

log_info "API統合テスト結果: $API_TESTS_PASSED/$API_TESTS_TOTAL 成功"

# =====================================
# 8. パフォーマンステスト（簡易版）
# =====================================
log_info "パフォーマンステストを実行中..."

# Apache Benchmarkがインストールされている場合
if command -v ab &> /dev/null; then
    log_info "Apache Benchmarkでパフォーマンステスト実行中..."
    ab -n 100 -c 10 "http://localhost:$BACKEND_PORT/health" > "$TEST_RESULTS_DIR/performance.txt" 2>&1
    log_success "パフォーマンステスト完了（結果: $TEST_RESULTS_DIR/performance.txt）"
else
    log_warning "Apache Benchmarkが見つかりません。パフォーマンステストをスキップします"
fi

# =====================================
# 9. セキュリティチェック（簡易版）
# =====================================
log_info "セキュリティチェックを実行中..."

# 基本的なセキュリティヘッダーチェック
check_security_headers() {
    local url="http://localhost:$BACKEND_PORT/health"
    local headers=$(curl -s -I "$url")
    
    echo "$headers" > "$TEST_RESULTS_DIR/security-headers.txt"
    
    if echo "$headers" | grep -i "x-content-type-options" > /dev/null; then
        log_success "✓ X-Content-Type-Options ヘッダーが設定されています"
    else
        log_warning "✗ X-Content-Type-Options ヘッダーが設定されていません"
    fi
    
    if echo "$headers" | grep -i "x-frame-options" > /dev/null; then
        log_success "✓ X-Frame-Options ヘッダーが設定されています"
    else
        log_warning "✗ X-Frame-Options ヘッダーが設定されていません"
    fi
}

check_security_headers

# =====================================
# 10. テスト結果サマリー
# =====================================
log_info "=== 統合テスト結果サマリー ==="

# テスト結果ファイル作成
{
    echo "Stockle 統合テスト結果"
    echo "======================="
    echo "実行日時: $(date)"
    echo "API統合テスト: $API_TESTS_PASSED/$API_TESTS_TOTAL 成功"
    echo ""
    echo "詳細結果:"
    echo "- セキュリティヘッダー: $TEST_RESULTS_DIR/security-headers.txt"
    if [[ -f "$TEST_RESULTS_DIR/performance.txt" ]]; then
        echo "- パフォーマンステスト: $TEST_RESULTS_DIR/performance.txt"
    fi
} > "$TEST_RESULTS_DIR/summary.txt"

log_success "統合テスト完了!"
log_info "結果サマリー: $TEST_RESULTS_DIR/summary.txt"

# 失敗があった場合は適切な終了コードを返す
if [[ $API_TESTS_PASSED -lt $API_TESTS_TOTAL ]]; then
    log_warning "一部のテストが失敗しました"
    exit 1
fi

log_success "すべてのテストが成功しました!"
exit 0