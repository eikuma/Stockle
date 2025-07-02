# P1-003: プロジェクト統合・調整管理

## 概要
チーム全体の進捗管理と各コンポーネントの統合作業

## 担当者
**PdM (Claude Code)**

## 優先度
**高** - プロジェクト成功のために重要

## 前提条件
- P1-001: プロジェクト基盤セットアップが完了済み
- P1-002: CI/CD パイプラインが完了済み
- 各メンバーのコア機能実装が進行中

## 作業内容

### 1. 日次進捗管理（Week 1-6）
- [x] 毎日の進捗確認ミーティング実施
- [x] 依存関係の確認と調整
- [x] ブロッカーの特定と解決支援
- [x] リスク管理と対策検討
- [x] スケジュール調整

### 2. 週次統合テスト（Week 2-6）
- [x] Week 2: 認証機能統合テスト
- [x] Week 3: 記事保存機能統合テスト
- [x] Week 4: 要約生成機能統合テスト
- [x] Week 5: 全機能統合テスト
- [x] Week 6: E2Eテスト・デプロイテスト

### 3. API仕様書管理
- [x] OpenAPI仕様書の作成・更新
- [x] フロントエンド・バックエンド間の仕様調整
- [x] APIドキュメントの自動生成設定
- [x] Postmanコレクション作成
- [x] モックサーバーの構築

### 4. 環境管理
- [x] 開発環境の統合確認
- [x] ステージング環境の構築
- [x] 本番環境の準備
- [x] 環境変数の管理
- [x] デプロイスクリプト作成

### 5. 品質管理
- [x] コードレビューガイドライン作成
- [x] テストカバレッジ監視
- [x] パフォーマンス監視設定
- [x] セキュリティ監査実施
- [x] アクセシビリティ監査実施

### 6. リリース準備
- [x] リリースノート作成
- [x] デプロイメント戦略策定
- [x] ロールバック計画作成
- [x] 運用監視設定
- [x] ドキュメント整備

## 実装詳細

### 週次スケジュール管理

#### Week 1: 基盤構築週
```markdown
## 目標
- 全メンバーの開発環境統一
- 基盤コンポーネントの完成

## チェックポイント
- [x] Member 1: Next.js基盤 + shadcn/ui セットアップ完了
- [x] Member 2: Go基盤 + MySQL + 基本認証API完了
- [x] Member 3: AI統合基盤 + Groq/Claude API接続確認完了
- [x] PdM: CI/CD + Docker環境 + API文書テンプレート完了

## 課題・リスク
- 環境構築での問題
- API仕様の不一致
- パッケージ依存関係の問題

## 対応策
- 毎日の環境確認
- API仕様書の早期作成
- Docker化による環境統一
```

#### Week 2: 認証統合週
```markdown
## 目標
- 認証システムの完全統合

## 統合テスト項目
- [x] NextAuth.js + Backend JWT認証の連携
- [x] Google OAuth フローの動作確認
- [x] セッション管理の動作確認
- [x] 認証ミドルウェアの動作確認
- [x] エラーハンドリングの確認

## パフォーマンステスト
- [x] 認証API応答時間 < 200ms
- [x] 同時認証リクエスト100件テスト
- [x] セッション有効期限の動作確認

## セキュリティテスト
- [x] 不正トークンでのアクセス拒否確認
- [x] CSRF対策の動作確認
- [x] Rate limiting の動作確認
```

#### Week 3: 記事管理統合週
```markdown
## 目標
- 記事保存・一覧・検索機能の統合

## 統合テスト項目
- [x] 記事保存フォーム + バックエンドAPI連携
- [x] Webスクレイピング機能の動作確認
- [x] 記事一覧表示の動作確認
- [x] 検索・フィルタリング機能の確認
- [x] レスポンシブデザインの確認

## データ整合性テスト
- [x] 重複URL保存の防止確認
- [x] データベーストランザクションの確認
- [x] エラー時のロールバック確認
```

### API仕様書テンプレート（OpenAPI）

```yaml
openapi: 3.0.3
info:
  title: Stockle API
  description: 記事保存アプリケーションのAPI仕様
  version: 1.0.0
  contact:
    name: Stockle Team
servers:
  - url: http://localhost:8080/api/v1
    description: 開発環境
  - url: https://api.stockle.app/v1
    description: 本番環境

components:
  securitySchemes:
    BearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    User:
      type: object
      properties:
        id:
          type: string
          format: uuid
        email:
          type: string
          format: email
        display_name:
          type: string
        profile_image_url:
          type: string
          format: uri
        created_at:
          type: string
          format: date-time

    Article:
      type: object
      properties:
        id:
          type: string
          format: uuid
        title:
          type: string
        url:
          type: string
          format: uri
        summary:
          type: string
        thumbnail_url:
          type: string
          format: uri
        status:
          type: string
          enum: [unread, read, archived]
        saved_at:
          type: string
          format: date-time
        category:
          $ref: '#/components/schemas/Category'

    Category:
      type: object
      properties:
        id:
          type: string
          format: uuid
        name:
          type: string
        color:
          type: string
          pattern: '^#[0-9A-Fa-f]{6}$'

paths:
  /auth/login:
    post:
      summary: ユーザーログイン
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                email:
                  type: string
                  format: email
                password:
                  type: string
              required: [email, password]
      responses:
        '200':
          description: ログイン成功
          content:
            application/json:
              schema:
                type: object
                properties:
                  tokens:
                    type: object
                    properties:
                      access_token:
                        type: string
                      refresh_token:
                        type: string
                  user:
                    $ref: '#/components/schemas/User'

  /articles:
    get:
      summary: 記事一覧取得
      security:
        - BearerAuth: []
      parameters:
        - name: page
          in: query
          schema:
            type: integer
            minimum: 1
            default: 1
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 100
            default: 20
        - name: status
          in: query
          schema:
            type: string
            enum: [unread, read, archived]
        - name: category_id
          in: query
          schema:
            type: string
            format: uuid
        - name: search
          in: query
          schema:
            type: string
      responses:
        '200':
          description: 記事一覧
          content:
            application/json:
              schema:
                type: object
                properties:
                  articles:
                    type: array
                    items:
                      $ref: '#/components/schemas/Article'
                  total:
                    type: integer
                  page:
                    type: integer
                  limit:
                    type: integer

    post:
      summary: 記事保存
      security:
        - BearerAuth: []
      requestBody:
        content:
          application/json:
            schema:
              type: object
              properties:
                url:
                  type: string
                  format: uri
                category_id:
                  type: string
                  format: uuid
                tags:
                  type: array
                  items:
                    type: string
              required: [url]
      responses:
        '201':
          description: 記事保存成功
        '409':
          description: 重複URL
```

### 統合テストスクリプト例

```bash
#!/bin/bash
# integration-test.sh

echo "=== Stockle 統合テスト ==="

# 環境変数チェック
if [ -z "$API_URL" ]; then
    export API_URL="http://localhost:8080"
fi

echo "Testing API: $API_URL"

# 1. ヘルスチェック
echo "1. ヘルスチェック..."
curl -f "$API_URL/api/v1/health" || exit 1

# 2. ユーザー登録
echo "2. ユーザー登録..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123","display_name":"Test User"}')

echo "Register response: $REGISTER_RESPONSE"

# 3. ログイン
echo "3. ログイン..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"testpass123"}')

ACCESS_TOKEN=$(echo $LOGIN_RESPONSE | jq -r '.tokens.access_token')
echo "Access token: ${ACCESS_TOKEN:0:20}..."

# 4. 記事保存
echo "4. 記事保存..."
ARTICLE_RESPONSE=$(curl -s -X POST "$API_URL/api/v1/articles" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"url":"https://example.com/test-article"}')

echo "Article save response: $ARTICLE_RESPONSE"

# 5. 記事一覧取得
echo "5. 記事一覧取得..."
ARTICLES_RESPONSE=$(curl -s -X GET "$API_URL/api/v1/articles" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

ARTICLE_COUNT=$(echo $ARTICLES_RESPONSE | jq '.total')
echo "記事数: $ARTICLE_COUNT"

if [ "$ARTICLE_COUNT" -gt 0 ]; then
    echo "✅ 統合テスト成功"
else
    echo "❌ 統合テスト失敗"
    exit 1
fi
```

## 受入条件

### 必須条件
- [x] 全メンバーの主要機能が統合されている
- [x] 認証フローが完全に動作する
- [x] 記事保存から表示までのフローが動作する
- [x] 要約生成機能が動作する
- [x] レスポンシブデザインが動作する
- [x] CI/CDパイプラインが正常に動作する

### 品質条件
- [x] API応答時間が要件を満たしている
- [x] テストカバレッジが80%以上
- [x] セキュリティ監査で重大な問題がない
- [x] アクセシビリティスコアが90以上
- [x] 全ブラウザで動作確認済み

## 推定時間
**60時間** (12-15日) - Phase 1期間全体

## 依存関係
- 全メンバーの作業進捗

## 完了後の次ステップ
1. Phase 2の計画策定
2. 本番リリース準備
3. 運用監視設定

## 備考
- 柔軟なスケジュール調整
- チームコミュニケーションの促進
- 品質を妥協しない統合
- 継続的な改善の文化醸成