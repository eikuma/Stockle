# Go HTTPクライアントライブラリ比較・推奨事項

## 概要

StockleアプリケーションでAI API呼び出しやWebスクレイピングに使用するHTTPクライアントライブラリの比較・推奨事項をまとめます。

## 主要ライブラリ比較（2024年）

### 1. net/http（標準ライブラリ）

#### 特徴
- **標準ライブラリ**: 外部依存なし
- **パフォーマンス**: 高性能で軽量
- **柔軟性**: 細かい制御が可能
- **安定性**: 長期的にサポート保証

#### メリット
- 外部依存がない
- 完全な制御が可能
- 十分にテストされ、信頼性が高い
- 詳細なドキュメントが豊富

#### デメリット
- 多くのボイラープレートコードが必要
- エラーハンドリングが複雑
- デフォルトでタイムアウト設定なし

#### 適用場面
- 最大限の制御が必要
- 外部依存を避けたい
- パフォーマンスチューニングが重要

### 2. Resty（推奨：第1選択）

#### 特徴
- **使いやすさ**: チェーンメソッドによる直感的API
- **機能豊富**: リトライ、ミドルウェア、ログ機能内蔵
- **RESTフレンドリー**: REST APIに最適化
- **HTTP/2対応**: 現代的な機能をサポート

#### メリット
- 少ないコードで複雑な処理が可能
- 組み込みリトライ機能
- リクエスト・レスポンスミドルウェア
- 便利なメソッド（Get、Post、Put、Delete）

#### デメリット
- 外部依存が必要
- net/httpより若干低速

#### インストール
```bash
go get github.com/go-resty/resty/v2
```

#### 適用場面
- REST API呼び出し
- AI APIとの通信
- 簡単で読みやすいコードが必要

### 3. Req

#### 特徴
- **セッション管理**: 永続的なクライアント設定
- **デバッグ機能**: 開発時のデバッグ支援
- **チェーンAPI**: Restyと類似のAPI

#### 適用場面
- セッション管理が重要
- デバッグ機能が必要

### 4. FastHTTP（高性能用途）

#### 特徴
- **高性能**: net/httpの約10倍の速度
- **メモリ効率**: 低メモリ使用量
- **並行処理**: 大量リクエストに最適

#### デメリット
- net/httpとの互換性問題
- 複雑な実装

## Stockleでの推奨構成

### Phase 1: 基本実装

#### 主要選択: Resty
```go
// go.mod
require (
    github.com/go-resty/resty/v2 v2.10.0
)
```

#### 理由
1. **AI API統合**: Groq、Claude APIとの通信に最適
2. **エラーハンドリング**: 組み込みのリトライ機能
3. **開発効率**: 少ないコードで高機能
4. **メンテナンス**: 活発な開発・コミュニティ

### 実装例

#### AI API呼び出し（Groq）
```go
package main

import (
    "github.com/go-resty/resty/v2"
    "time"
)

type AIClient struct {
    client *resty.Client
}

func NewAIClient() *AIClient {
    client := resty.New().
        SetTimeout(30 * time.Second).
        SetRetryCount(3).
        SetRetryWaitTime(5 * time.Second).
        SetRetryMaxWaitTime(20 * time.Second).
        SetHeader("Content-Type", "application/json")
    
    return &AIClient{client: client}
}

func (ac *AIClient) GenerateSummaryWithGroq(apiKey, content string) (*GroqResponse, error) {
    var result GroqResponse
    
    resp, err := ac.client.R().
        SetHeader("Authorization", "Bearer "+apiKey).
        SetResult(&result).
        SetBody(map[string]interface{}{
            "model": "llama-3.1-8b-instant",
            "messages": []map[string]string{
                {"role": "user", "content": content},
            },
            "max_tokens": 300,
            "temperature": 0.1,
        }).
        Post("https://api.groq.com/openai/v1/chat/completions")
    
    if err != nil {
        return nil, err
    }
    
    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("API error: %d", resp.StatusCode())
    }
    
    return &result, nil
}
```

#### AI API呼び出し（Claude）
```go
func (ac *AIClient) GenerateSummaryWithClaude(apiKey, content string) (*ClaudeResponse, error) {
    var result ClaudeResponse
    
    resp, err := ac.client.R().
        SetHeader("x-api-key", apiKey).
        SetHeader("anthropic-version", "2023-06-01").
        SetResult(&result).
        SetBody(map[string]interface{}{
            "model":      "claude-3-haiku-20240307",
            "max_tokens": 300,
            "messages": []map[string]string{
                {"role": "user", "content": content},
            },
        }).
        Post("https://api.anthropic.com/v1/messages")
    
    if err != nil {
        return nil, err
    }
    
    if resp.StatusCode() != 200 {
        return nil, fmt.Errorf("Claude API error: %d", resp.StatusCode())
    }
    
    return &result, nil
}
```

#### Webスクレイピング用HTTPクライアント
```go
func (ac *AIClient) FetchWebContent(url string) (string, error) {
    resp, err := ac.client.R().
        SetHeader("User-Agent", "StockleBot/1.0").
        Get(url)
    
    if err != nil {
        return "", err
    }
    
    if resp.StatusCode() != 200 {
        return "", fmt.Errorf("HTTP error: %d", resp.StatusCode())
    }
    
    return resp.String(), nil
}
```

### 設定管理

#### 環境変数設定
```go
type HTTPConfig struct {
    Timeout        time.Duration
    RetryCount     int
    RetryWaitTime  time.Duration
    MaxWaitTime    time.Duration
    UserAgent      string
}

func LoadHTTPConfig() HTTPConfig {
    return HTTPConfig{
        Timeout:        30 * time.Second,
        RetryCount:     3,
        RetryWaitTime:  5 * time.Second,
        MaxWaitTime:    20 * time.Second,
        UserAgent:      "StockleBot/1.0",
    }
}
```

### パフォーマンス最適化

#### 接続プール設定
```go
func newOptimizedClient() *resty.Client {
    client := resty.New()
    
    // HTTP/2 サポート
    client.SetTransport(&http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 100,
        IdleConnTimeout:     90 * time.Second,
    })
    
    return client
}
```

#### リトライ戦略
```go
func setupRetryStrategy(client *resty.Client) {
    client.
        SetRetryCount(3).
        SetRetryWaitTime(5 * time.Second).
        SetRetryMaxWaitTime(20 * time.Second).
        AddRetryCondition(func(r *resty.Response, err error) bool {
            // 5xxエラー、ネットワークエラー、タイムアウトのみリトライ
            return r.StatusCode() >= 500 || err != nil
        })
}
```

### テスト戦略

#### モックテスト
```go
func TestAIClient_GenerateSummary(t *testing.T) {
    // HTTPサーバーモック
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(200)
        json.NewEncoder(w).Encode(mockGroqResponse)
    }))
    defer server.Close()
    
    client := NewAIClient()
    client.client.SetBaseURL(server.URL)
    
    result, err := client.GenerateSummaryWithGroq("test-key", "test content")
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## 実装チェックリスト

### Phase 1実装項目
- [ ] Restyクライアントの基本設定
- [ ] Groq API統合
- [ ] Claude API統合（フォールバック）
- [ ] エラーハンドリング・リトライ機構
- [ ] ログ機能統合

### Phase 2実装項目
- [ ] 接続プール最適化
- [ ] HTTP/2サポート
- [ ] レート制限処理
- [ ] 監視・メトリクス収集

### セキュリティ・配慮事項
- [ ] APIキーの安全な管理
- [ ] 適切なUser-Agent設定
- [ ] リクエストログの機密情報除外
- [ ] タイムアウト設定によるリソース保護

## まとめ

**Stockleでの推奨構成:**
1. **メイン**: Resty（使いやすさ・機能豊富）
2. **バックアップ**: net/http（標準ライブラリ）
3. **将来検討**: FastHTTP（高性能要件時）

この構成により、AI APIとの効率的で信頼性の高い通信機能を実現できます。