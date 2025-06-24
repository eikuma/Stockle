# ジョブキュー機能設計書

## 概要

Stockleアプリケーションにおける非同期処理（AI要約生成、音声ファイル生成等）のためのジョブキュー機能の設計書です。

## アーキテクチャ

### 1. ジョブキュー実装方針

#### 1.1 データベースベースのジョブキュー
- **主要ストア**: MySQL（job_queueテーブル）
- **キャッシュ**: Redis（高速な状態管理）
- **ワーカー**: Goルーチンによる並行処理

#### 1.2 ハイブリッド構成
```
[API Request] → [Job Creation] → [MySQL job_queue] → [Redis Notification] → [Worker Pool]
```

### 2. ジョブタイプ定義

#### 2.1 Phase 1 ジョブタイプ
- `summarize`: AI要約生成
- `extract_content`: 記事本文抽出
- `validate_summary`: 要約品質チェック

#### 2.2 Phase 2 ジョブタイプ（今後追加予定）
- `generate_podcast`: 音声ファイル生成
- `calculate_similarity`: 記事類似度計算
- `batch_process`: バッチ処理

### 3. ジョブ優先度システム

#### 3.1 優先度レベル
1. **CRITICAL (10)**: システムクリティカルな処理
2. **HIGH (7-9)**: ユーザー待機中の処理
3. **NORMAL (4-6)**: 通常の非同期処理
4. **LOW (1-3)**: バックグラウンド処理

#### 3.2 優先度決定ロジック
```go
func calculateJobPriority(jobType string, userTier string) int {
    switch jobType {
    case "summarize":
        return 7 // HIGH
    case "generate_podcast":
        return 5 // NORMAL
    case "calculate_similarity":
        return 3 // LOW
    default:
        return 5 // NORMAL
    }
}
```

### 4. リトライ機構

#### 4.1 リトライ戦略
- **最大リトライ回数**: 3回
- **バックオフ戦略**: 指数バックオフ（2^n秒）
- **リトライ間隔**: 1秒 → 2秒 → 4秒 → 失敗

#### 4.2 エラー分類
```go
type ErrorType string

const (
    RetryableError    ErrorType = "retryable"    // API timeout, network error
    NonRetryableError ErrorType = "non_retryable" // auth error, invalid request
    RateLimitError    ErrorType = "rate_limit"    // API rate limit exceeded
)
```

### 5. ワーカー実装

#### 5.1 ワーカープール設計
```go
type WorkerPool struct {
    workerCount    int
    jobQueue       chan Job
    quit           chan bool
    wg             *sync.WaitGroup
    rateLimiter    *time.Ticker
}

type Job struct {
    ID       string
    Type     string
    Priority int
    Payload  map[string]interface{}
    Retries  int
    MaxRetries int
}
```

#### 5.2 ワーカー数の動的調整
- **CPU使用率**: 70%以下を維持
- **メモリ使用率**: 80%以下を維持
- **動的スケーリング**: 負荷に応じて1-10ワーカーで調整

### 6. 監視・ログ機能

#### 6.1 メトリクス収集
- ジョブ実行時間
- 成功・失敗率
- キュー長
- ワーカー使用率

#### 6.2 ログレベル
```go
type LogLevel string

const (
    DEBUG LogLevel = "DEBUG"
    INFO  LogLevel = "INFO"
    WARN  LogLevel = "WARN"
    ERROR LogLevel = "ERROR"
)
```

### 7. API設計

#### 7.1 ジョブ作成API
```http
POST /api/v1/jobs
Content-Type: application/json

{
    "type": "summarize",
    "payload": {
        "article_id": "uuid",
        "user_preferences": {
            "length": "medium",
            "language": "ja"
        }
    },
    "priority": 7
}
```

#### 7.2 ジョブ状態確認API
```http
GET /api/v1/jobs/{job_id}

Response:
{
    "job_id": "uuid",
    "status": "completed",
    "progress": 100,
    "result": {
        "summary": "生成された要約文...",
        "generated_at": "2024-01-01T00:00:00Z"
    },
    "error": null
}
```

### 8. AI API統合

#### 8.1 AI プロバイダー管理
```go
type AIProvider interface {
    GenerateSummary(content string) (*SummaryResult, error)
    GetRateLimit() RateLimit
    GetCost() float64
}

type GroqProvider struct {
    apiKey string
    client *http.Client
}

type ClaudeProvider struct {
    apiKey string
    client *http.Client
}
```

#### 8.2 フォールバック戦略
1. **Groq API**: 第一選択（無料枠）
2. **Claude API**: 第二選択（従量課金）
3. **Simple Extraction**: 最終手段（記事冒頭300文字）

### 9. 設定管理

#### 9.1 環境変数
```env
# Job Queue Configuration
JOB_QUEUE_WORKER_COUNT=5
JOB_QUEUE_MAX_RETRIES=3
JOB_QUEUE_BATCH_SIZE=10
JOB_QUEUE_POLL_INTERVAL=5s

# AI API Configuration
AI_REQUEST_TIMEOUT=30s
AI_MAX_CONCURRENT_REQUESTS=10
AI_FALLBACK_ENABLED=true

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=100
RATE_LIMIT_BURST_SIZE=10
```

### 10. 実装フェーズ

#### 10.1 Phase 1: 基本実装
- [ ] データベーステーブル作成
- [ ] 基本ワーカープール実装
- [ ] Summarizeジョブ実装
- [ ] Groq API統合

#### 10.2 Phase 2: 拡張機能
- [ ] Redis統合
- [ ] Claude API統合
- [ ] 動的スケーリング
- [ ] 詳細監視

#### 10.3 Phase 3: 最適化
- [ ] パフォーマンス最適化
- [ ] 高可用性対応
- [ ] 負荷分散

### 11. テスト戦略

#### 11.1 単体テスト
- ジョブ作成・実行・リトライ
- エラーハンドリング
- 優先度処理

#### 11.2 統合テスト
- API統合テスト
- データベース整合性
- リアルタイム通知

### 12. パフォーマンス目標

#### 12.1 レスポンス時間
- ジョブ作成: 50ms以内
- 要約生成: 30秒以内
- 状態確認: 10ms以内

#### 12.2 スループット
- 最大同時ジョブ数: 100
- ジョブ処理能力: 1000件/時

この設計により、スケーラブルで信頼性の高いジョブキュー機能を実現します。