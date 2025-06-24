# 記事保存アプリ - API設計書

## 1. 概要

### 1.1 API設計方針

- RESTfulアーキテクチャに準拠
- JSON形式でのデータ交換
- UTF-8エンコーディング
- セマンティックなHTTPステータスコードの使用
- 一貫性のあるエラーレスポンス形式

### 1.2 ベースURL

```
https://api.readlater.app/v1
```

### 1.3 認証方式

- JWT (JSON Web Token) を使用
- Authorization ヘッダーに Bearer トークンとして設定
- トークン有効期限: 7日間（リフレッシュトークン: 30日間）

```
Authorization: Bearer <jwt_token>
```

### 1.4 共通ヘッダー

**リクエストヘッダー**
```
Content-Type: application/json
Accept: application/json
Authorization: Bearer <jwt_token>
```

**レスポンスヘッダー**
```
Content-Type: application/json; charset=utf-8
X-Request-ID: <unique_request_id>
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 99
X-RateLimit-Reset: 1640995200
```

### 1.5 エラーレスポンス形式

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "入力値が不正です",
    "details": [
      {
        "field": "email",
        "message": "有効なメールアドレスを入力してください"
      }
    ]
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

## 2. 認証関連エンドポイント

### 2.1 ユーザー登録

**エンドポイント**
```
POST /auth/register
```

**リクエストボディ**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "display_name": "山田太郎"
}
```

**レスポンス**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "display_name": "山田太郎",
    "created_at": "2024-01-01T00:00:00Z"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 604800
  }
}
```

**ステータスコード**
- 201: 作成成功
- 400: バリデーションエラー
- 409: メールアドレス重複

### 2.2 ログイン

**エンドポイント**
```
POST /auth/login
```

**リクエストボディ**
```json
{
  "email": "user@example.com",
  "password": "securePassword123",
  "remember_me": true
}
```

**レスポンス**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "display_name": "山田太郎",
    "profile_image_url": "https://example.com/avatar.jpg"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 604800
  }
}
```

### 2.3 Google OAuth認証

**認証開始**
```
GET /auth/google
```

**コールバック**
```
GET /auth/google/callback?code=<authorization_code>&state=<state>
```

**レスポンス**
```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "display_name": "山田太郎",
    "profile_image_url": "https://lh3.googleusercontent.com/...",
    "auth_provider": "google"
  },
  "tokens": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "expires_in": 604800
  }
}
```

### 2.4 トークンリフレッシュ

**エンドポイント**
```
POST /auth/refresh
```

**リクエストボディ**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### 2.5 ログアウト

**エンドポイント**
```
POST /auth/logout
```

**リクエストボディ**
```json
{
  "all_devices": false
}
```

## 3. 記事管理エンドポイント

### 3.1 記事保存

**エンドポイント**
```
POST /articles
```

**リクエストボディ**
```json
{
  "url": "https://example.com/article",
  "category_id": "550e8400-e29b-41d4-a716-446655440000",
  "tags": ["技術", "AI"]
}
```

**レスポンス**
```json
{
  "article": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "url": "https://example.com/article",
    "title": "AIの未来について",
    "summary": "人工知能技術の発展により...",
    "content_snippet": "本記事では、最新のAI技術...",
    "thumbnail_url": "https://example.com/thumb.jpg",
    "author": "技術太郎",
    "published_at": "2024-01-01T00:00:00Z",
    "saved_at": "2024-01-02T00:00:00Z",
    "category": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "テクノロジー"
    },
    "tags": ["技術", "AI"],
    "status": "unread",
    "summary_generation_status": "completed"
  }
}
```

**非同期処理**
- 要約生成は非同期で実行
- WebSocketまたはSSEで進捗通知

### 3.2 記事一覧取得

**エンドポイント**
```
GET /articles
```

**クエリパラメータ**
```
?page=1
&per_page=20
&sort=saved_at
&order=desc
&status=unread
&category_id=550e8400-e29b-41d4-a716-446655440000
&tag=AI
&search=人工知能
&saved_from=2024-01-01
&saved_to=2024-01-31
```

**レスポンス**
```json
{
  "articles": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "url": "https://example.com/article",
      "title": "AIの未来について",
      "summary": "人工知能技術の発展により...",
      "thumbnail_url": "https://example.com/thumb.jpg",
      "saved_at": "2024-01-02T00:00:00Z",
      "category": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "テクノロジー"
      },
      "tags": ["技術", "AI"],
      "status": "unread",
      "reading_time_minutes": 5
    }
  ],
  "pagination": {
    "current_page": 1,
    "per_page": 20,
    "total_pages": 5,
    "total_items": 98,
    "has_next": true,
    "has_previous": false
  }
}
```

### 3.3 記事詳細取得

**エンドポイント**
```
GET /articles/{article_id}
```

**レスポンス**
```json
{
  "article": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "url": "https://example.com/article",
    "title": "AIの未来について",
    "summary": "人工知能技術の発展により...",
    "content": "記事の本文全体...",
    "thumbnail_url": "https://example.com/thumb.jpg",
    "author": "技術太郎",
    "published_at": "2024-01-01T00:00:00Z",
    "saved_at": "2024-01-02T00:00:00Z",
    "last_accessed_at": "2024-01-03T00:00:00Z",
    "category": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "テクノロジー",
      "confidence_score": 0.85
    },
    "tags": ["技術", "AI"],
    "status": "read",
    "reading_progress": 0.75,
    "metadata": {
      "site_name": "Tech Blog",
      "description": "最新の技術トレンドを紹介",
      "og_image": "https://example.com/og.jpg"
    }
  }
}
```

### 3.4 記事更新

**エンドポイント**
```
PATCH /articles/{article_id}
```

**リクエストボディ**
```json
{
  "category_id": "550e8400-e29b-41d4-a716-446655440000",
  "tags": ["技術", "AI", "機械学習"],
  "status": "read"
}
```

### 3.5 記事削除

**エンドポイント**
```
DELETE /articles/{article_id}
```

### 3.6 記事一括操作

**エンドポイント**
```
POST /articles/bulk
```

**リクエストボディ**
```json
{
  "article_ids": [
    "123e4567-e89b-12d3-a456-426614174000",
    "223e4567-e89b-12d3-a456-426614174001"
  ],
  "operation": "mark_as_read"
}
```

**操作タイプ**
- mark_as_read: 既読にする
- mark_as_unread: 未読にする
- delete: 削除
- change_category: カテゴリ変更
- add_tags: タグ追加
- remove_tags: タグ削除

## 4. カテゴリ管理エンドポイント

### 4.1 カテゴリ一覧取得

**エンドポイント**
```
GET /categories
```

**レスポンス**
```json
{
  "categories": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "テクノロジー",
      "color": "#3B82F6",
      "article_count": 42,
      "is_default": true,
      "created_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### 4.2 カテゴリ作成

**エンドポイント**
```
POST /categories
```

**リクエストボディ**
```json
{
  "name": "プログラミング",
  "color": "#10B981"
}
```

### 4.3 カテゴリ更新

**エンドポイント**
```
PUT /categories/{category_id}
```

### 4.4 カテゴリ削除

**エンドポイント**
```
DELETE /categories/{category_id}
```

## 5. 要約・AI関連エンドポイント

### 5.1 要約再生成

**エンドポイント**
```
POST /articles/{article_id}/summarize
```

**リクエストボディ**
```json
{
  "length": "medium",
  "language": "ja"
}
```

**レスポンス**
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "processing",
  "estimated_time_seconds": 30
}
```

### 5.2 要約生成状態確認

**エンドポイント**
```
GET /jobs/{job_id}
```

**レスポンス**
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "completed",
  "result": {
    "summary": "新しく生成された要約...",
    "generated_at": "2024-01-01T00:00:00Z"
  }
}
```

## 6. 検索エンドポイント

### 6.1 全文検索

**エンドポイント**
```
GET /search
```

**クエリパラメータ**
```
?q=人工知能
&type=all
&category_id=550e8400-e29b-41d4-a716-446655440000
&status=unread
&from=2024-01-01
&to=2024-01-31
&page=1
&per_page=20
```

**レスポンス**
```json
{
  "results": [
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "title": "AIの未来について",
      "summary": "人工知能技術の発展により...",
      "url": "https://example.com/article",
      "matched_fields": ["title", "summary"],
      "highlight": {
        "title": "<mark>AI</mark>の未来について",
        "summary": "<mark>人工知能</mark>技術の発展により..."
      }
    }
  ],
  "total_results": 15,
  "search_time_ms": 120
}
```

## 7. ユーザー設定エンドポイント

### 7.1 ユーザー情報取得

**エンドポイント**
```
GET /users/me
```

### 7.2 ユーザー情報更新

**エンドポイント**
```
PATCH /users/me
```

**リクエストボディ**
```json
{
  "display_name": "新しい名前",
  "default_category_id": "550e8400-e29b-41d4-a716-446655440000",
  "preferences": {
    "items_per_page": 30,
    "default_sort": "saved_at",
    "default_view": "grid"
  }
}
```

### 7.3 パスワード変更

**エンドポイント**
```
POST /users/me/change-password
```

### 7.4 アカウント削除

**エンドポイント**
```
DELETE /users/me
```

## 8. Phase 2: 音声関連エンドポイント

### 8.1 ポッドキャスト生成

**エンドポイント**
```
POST /articles/{article_id}/podcast
```

**リクエストボディ**
```json
{
  "voice_type": "female",
  "speed": 1.0,
  "language": "ja"
}
```

### 8.2 ポッドキャスト取得

**エンドポイント**
```
GET /articles/{article_id}/podcast
```

**レスポンス**
```json
{
  "podcast": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "article_id": "123e4567-e89b-12d3-a456-426614174000",
    "audio_url": "https://cdn.example.com/podcasts/xxx.mp3",
    "duration_seconds": 180,
    "file_size_bytes": 2880000,
    "voice_type": "female",
    "speed": 1.0,
    "generated_at": "2024-01-01T00:00:00Z"
  }
}
```

### 8.3 プレイリスト作成

**エンドポイント**
```
POST /playlists
```

### 8.4 再生履歴記録

**エンドポイント**
```
POST /podcasts/{podcast_id}/history
```

**リクエストボディ**
```json
{
  "position_seconds": 120,
  "completed": false
}
```

## 9. Phase 3: 統計・分析エンドポイント

### 9.1 読書統計取得

**エンドポイント**
```
GET /stats/reading
```

**クエリパラメータ**
```
?period=week
&from=2024-01-01
&to=2024-01-07
```

**レスポンス**
```json
{
  "stats": {
    "total_reading_time_minutes": 245,
    "articles_read": 12,
    "articles_saved": 20,
    "average_reading_time_minutes": 20.4,
    "categories_breakdown": [
      {
        "category": "テクノロジー",
        "count": 8,
        "percentage": 66.7
      }
    ],
    "daily_stats": [
      {
        "date": "2024-01-01",
        "reading_time_minutes": 45,
        "articles_read": 3
      }
    ]
  }
}
```

### 9.2 類似記事取得

**エンドポイント**
```
GET /articles/{article_id}/similar
```

## 10. Phase 4: ダイジェスト・レコメンデーション

### 10.1 ダイジェスト設定

**エンドポイント**
```
POST /digest/settings
```

**リクエストボディ**
```json
{
  "enabled": true,
  "frequency": "daily",
  "time": "09:00",
  "timezone": "Asia/Tokyo",
  "categories": ["550e8400-e29b-41d4-a716-446655440000"],
  "max_articles": 10,
  "include_stats": true
}
```

### 10.2 レコメンデーション取得

**エンドポイント**
```
GET /recommendations
```

**レスポンス**
```json
{
  "recommendations": [
    {
      "article": {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "title": "おすすめの記事",
        "summary": "あなたの興味に基づいて..."
      },
      "score": 0.92,
      "reason": "「AI」タグの記事をよく読んでいるため"
    }
  ]
}
```

## 11. WebSocket エンドポイント

### 11.1 リアルタイム通知

**エンドポイント**
```
wss://api.readlater.app/v1/ws
```

**イベントタイプ**
```json
{
  "type": "summary_completed",
  "data": {
    "article_id": "123e4567-e89b-12d3-a456-426614174000",
    "summary": "生成された要約..."
  }
}
```

## 12. レート制限

### 12.1 制限値

- 認証済みユーザー: 100リクエスト/分
- 未認証ユーザー: 10リクエスト/分
- 記事保存: 20記事/時間

### 12.2 制限超過時のレスポンス

```json
{
  "error": {
    "code": "RATE_LIMIT_EXCEEDED",
    "message": "リクエスト数が制限を超えました",
    "retry_after_seconds": 60
  }
}
```

**ステータスコード**: 429 Too Many Requests

## 13. ステータスコード一覧

- 200: OK - 成功
- 201: Created - リソース作成成功
- 204: No Content - 成功（レスポンスボディなし）
- 400: Bad Request - リクエスト不正
- 401: Unauthorized - 認証エラー
- 403: Forbidden - アクセス権限なし
- 404: Not Found - リソースが見つからない
- 409: Conflict - リソースの競合
- 422: Unprocessable Entity - バリデーションエラー
- 429: Too Many Requests - レート制限超過
- 500: Internal Server Error - サーバーエラー
- 503: Service Unavailable - メンテナンス中