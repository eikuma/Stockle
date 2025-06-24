# P1-301: AI統合基盤構築

## 概要
Groq API + Anthropic Claude APIによる要約生成システムの構築

## 担当者
**Member 3 (Backend Features Developer)**

## 優先度
**高** - コア機能の実装

## 前提条件
- P1-201: バックエンド基盤構築が完了済み
- Groq API + Anthropic APIキーが取得済み
- 記事保存機能の基本実装が完了済み

## 作業内容

### 1. AI統合基盤の設計
- [ ] AI プロバイダーの抽象化インターフェース
- [ ] フォールバック機構の設計
- [ ] 非同期処理システムの設計
- [ ] エラーハンドリング戦略
- [ ] レート制限管理

### 2. Groq API統合
- [ ] Groq APIクライアントの実装
- [ ] API認証の実装
- [ ] 要約生成リクエストの実装
- [ ] レスポンス処理
- [ ] エラーハンドリング

### 3. Anthropic Claude API統合
- [ ] Claude APIクライアントの実装
- [ ] API認証の実装
- [ ] 要約生成リクエストの実装
- [ ] レスポンス処理
- [ ] フォールバック機能

### 4. 非同期ジョブシステム
- [ ] ジョブキューの実装
- [ ] ワーカープールの実装
- [ ] ジョブリトライ機構
- [ ] ジョブ状態管理
- [ ] 優先度制御

### 5. 要約生成サービス
- [ ] 記事コンテンツの前処理
- [ ] プロンプトエンジニアリング
- [ ] 要約品質の評価
- [ ] 複数バージョンの要約生成
- [ ] 結果の後処理

### 6. キャッシュシステム
- [ ] 要約結果のキャッシュ
- [ ] 重複処理の防止
- [ ] キャッシュ無効化
- [ ] 容量管理

## 実装詳細

### internal/services/ai_service.go
```go
package services

import (
    "context"
    "fmt"
    "strings"
    "time"
    
    "github.com/eikuma/stockle/backend/internal/config"
    "github.com/eikuma/stockle/backend/internal/models"
    "github.com/eikuma/stockle/backend/pkg/groq"
    "github.com/eikuma/stockle/backend/pkg/anthropic"
)

type AIService struct {
    groqClient      *groq.Client
    anthropicClient *anthropic.Client
    config          *config.AIConfig
}

type SummaryRequest struct {
    Content     string
    Title       string
    URL         string
    Language    string
    SummaryType string // "short", "medium", "long"
}

type SummaryResponse struct {
    Summary      string
    Confidence   float64
    Provider     string
    GeneratedAt  time.Time
    ModelVersion string
    WordCount    int
}

func NewAIService(cfg *config.AIConfig) *AIService {
    return &AIService{
        groqClient:      groq.NewClient(cfg.GroqAPIKey),
        anthropicClient: anthropic.NewClient(cfg.AnthropicAPIKey),
        config:          cfg,
    }
}

func (s *AIService) GenerateSummary(ctx context.Context, req *SummaryRequest) (*SummaryResponse, error) {
    // Groqを第一選択として使用
    summary, err := s.generateWithGroq(ctx, req)
    if err != nil {
        // Groqが失敗した場合はClaude APIにフォールバック
        summary, err = s.generateWithClaude(ctx, req)
        if err != nil {
            return nil, fmt.Errorf("all AI providers failed: %w", err)
        }
    }
    
    return summary, nil
}

func (s *AIService) generateWithGroq(ctx context.Context, req *SummaryRequest) (*SummaryResponse, error) {
    prompt := s.buildPrompt(req)
    
    groqReq := &groq.ChatCompletionRequest{
        Model: "llama3-8b-8192",
        Messages: []groq.Message{
            {
                Role:    "system",
                Content: s.getSystemPrompt(req.SummaryType),
            },
            {
                Role:    "user",
                Content: prompt,
            },
        },
        MaxTokens:   500,
        Temperature: 0.3,
    }
    
    resp, err := s.groqClient.CreateChatCompletion(ctx, groqReq)
    if err != nil {
        return nil, fmt.Errorf("groq API error: %w", err)
    }
    
    summary := strings.TrimSpace(resp.Choices[0].Message.Content)
    confidence := s.calculateConfidence(summary, req.Content)
    
    return &SummaryResponse{
        Summary:      summary,
        Confidence:   confidence,
        Provider:     "groq",
        GeneratedAt:  time.Now(),
        ModelVersion: resp.Model,
        WordCount:    len(strings.Fields(summary)),
    }, nil
}

func (s *AIService) generateWithClaude(ctx context.Context, req *SummaryRequest) (*SummaryResponse, error) {
    prompt := s.buildPrompt(req)
    
    claudeReq := &anthropic.MessageRequest{
        Model:     "claude-3-haiku-20240307",
        MaxTokens: 500,
        Messages: []anthropic.Message{
            {
                Role:    "user",
                Content: prompt,
            },
        },
        System: s.getSystemPrompt(req.SummaryType),
    }
    
    resp, err := s.anthropicClient.CreateMessage(ctx, claudeReq)
    if err != nil {
        return nil, fmt.Errorf("claude API error: %w", err)
    }
    
    summary := strings.TrimSpace(resp.Content[0].Text)
    confidence := s.calculateConfidence(summary, req.Content)
    
    return &SummaryResponse{
        Summary:      summary,
        Confidence:   confidence,
        Provider:     "claude",
        GeneratedAt:  time.Now(),
        ModelVersion: resp.Model,
        WordCount:    len(strings.Fields(summary)),
    }, nil
}

func (s *AIService) buildPrompt(req *SummaryRequest) string {
    return fmt.Sprintf(`記事タイトル: %s

記事URL: %s

記事内容:
%s

上記の記事を%sで要約してください。`, req.Title, req.URL, req.Content, s.getSummaryTypeDescription(req.SummaryType))
}

func (s *AIService) getSystemPrompt(summaryType string) string {
    switch summaryType {
    case "short":
        return "あなたは記事の要約を作成する専門家です。記事の主要なポイントを50-100文字で簡潔に要約してください。読者が記事の概要を素早く理解できるようにしてください。"
    case "long":
        return "あなたは記事の要約を作成する専門家です。記事の詳細な内容を500-800文字で包括的に要約してください。重要なポイント、詳細、結論を含めて、読者が記事を読まなくても内容を十分理解できるようにしてください。"
    default: // medium
        return "あなたは記事の要約を作成する専門家です。記事の主要なポイントを200-300文字で要約してください。重要な情報を含みつつ、読みやすい長さにまとめてください。"
    }
}

func (s *AIService) getSummaryTypeDescription(summaryType string) string {
    switch summaryType {
    case "short":
        return "50-100文字の短い要約"
    case "long":
        return "500-800文字の詳細な要約"
    default:
        return "200-300文字の標準的な要約"
    }
}

func (s *AIService) calculateConfidence(summary, originalContent string) float64 {
    // 簡単な信頼度計算（実際にはより複雑なロジックが必要）
    if len(summary) < 50 {
        return 0.6
    }
    if len(summary) > len(originalContent) {
        return 0.4
    }
    return 0.8
}
```

### pkg/groq/client.go
```go
package groq

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

const (
    BaseURL = "https://api.groq.com/openai/v1"
)

type Client struct {
    httpClient *http.Client
    apiKey     string
    baseURL    string
}

type ChatCompletionRequest struct {
    Model       string    `json:"model"`
    Messages    []Message `json:"messages"`
    MaxTokens   int       `json:"max_tokens,omitempty"`
    Temperature float64   `json:"temperature,omitempty"`
    TopP        float64   `json:"top_p,omitempty"`
    Stream      bool      `json:"stream,omitempty"`
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type ChatCompletionResponse struct {
    ID      string   `json:"id"`
    Object  string   `json:"object"`
    Created int64    `json:"created"`
    Model   string   `json:"model"`
    Choices []Choice `json:"choices"`
    Usage   Usage    `json:"usage"`
}

type Choice struct {
    Index        int     `json:"index"`
    Message      Message `json:"message"`
    FinishReason string  `json:"finish_reason"`
}

type Usage struct {
    PromptTokens     int `json:"prompt_tokens"`
    CompletionTokens int `json:"completion_tokens"`
    TotalTokens      int `json:"total_tokens"`
}

func NewClient(apiKey string) *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        apiKey:  apiKey,
        baseURL: BaseURL,
    }
}

func (c *Client) CreateChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
    
    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
    }
    
    var chatResp ChatCompletionResponse
    if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &chatResp, nil
}
```

### pkg/anthropic/client.go
```go
package anthropic

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

const (
    BaseURL = "https://api.anthropic.com/v1"
)

type Client struct {
    httpClient *http.Client
    apiKey     string
    baseURL    string
}

type MessageRequest struct {
    Model     string    `json:"model"`
    MaxTokens int       `json:"max_tokens"`
    Messages  []Message `json:"messages"`
    System    string    `json:"system,omitempty"`
}

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type MessageResponse struct {
    ID      string    `json:"id"`
    Type    string    `json:"type"`
    Role    string    `json:"role"`
    Content []Content `json:"content"`
    Model   string    `json:"model"`
    Usage   Usage     `json:"usage"`
}

type Content struct {
    Type string `json:"type"`
    Text string `json:"text"`
}

type Usage struct {
    InputTokens  int `json:"input_tokens"`
    OutputTokens int `json:"output_tokens"`
}

func NewClient(apiKey string) *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        apiKey:  apiKey,
        baseURL: BaseURL,
    }
}

func (c *Client) CreateMessage(ctx context.Context, req *MessageRequest) (*MessageResponse, error) {
    jsonData, err := json.Marshal(req)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/messages", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("x-api-key", c.apiKey)
    httpReq.Header.Set("anthropic-version", "2023-06-01")
    
    resp, err := c.httpClient.Do(httpReq)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
    }
    
    var msgResp MessageResponse
    if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
        return nil, fmt.Errorf("failed to decode response: %w", err)
    }
    
    return &msgResp, nil
}
```

### internal/services/job_service.go
```go
package services

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"
    
    "github.com/eikuma/stockle/backend/internal/models"
    "github.com/eikuma/stockle/backend/internal/repositories"
)

type JobService struct {
    jobRepo     repositories.JobRepository
    articleRepo repositories.ArticleRepository
    aiService   *AIService
}

type JobPayload struct {
    ArticleID string                 `json:"article_id"`
    JobType   string                 `json:"job_type"`
    Options   map[string]interface{} `json:"options"`
}

func NewJobService(jobRepo repositories.JobRepository, articleRepo repositories.ArticleRepository, aiService *AIService) *JobService {
    return &JobService{
        jobRepo:     jobRepo,
        articleRepo: articleRepo,
        aiService:   aiService,
    }
}

func (s *JobService) EnqueueSummaryJob(articleID string, priority int) error {
    payload := JobPayload{
        ArticleID: articleID,
        JobType:   "summarize",
        Options: map[string]interface{}{
            "summary_type": "medium",
        },
    }
    
    payloadJSON, err := json.Marshal(payload)
    if err != nil {
        return fmt.Errorf("failed to marshal job payload: %w", err)
    }
    
    job := &models.JobQueue{
        JobType:    "summarize",
        Priority:   priority,
        Status:     "pending",
        Payload:    string(payloadJSON),
        MaxRetries: 3,
    }
    
    return s.jobRepo.Create(job)
}

func (s *JobService) ProcessJob(ctx context.Context, job *models.JobQueue) error {
    var payload JobPayload
    if err := json.Unmarshal([]byte(job.Payload), &payload); err != nil {
        return fmt.Errorf("failed to unmarshal job payload: %w", err)
    }
    
    switch job.JobType {
    case "summarize":
        return s.processSummaryJob(ctx, job, &payload)
    default:
        return fmt.Errorf("unknown job type: %s", job.JobType)
    }
}

func (s *JobService) processSummaryJob(ctx context.Context, job *models.JobQueue, payload *JobPayload) error {
    // 記事の取得
    article, err := s.articleRepo.GetByID(payload.ArticleID)
    if err != nil {
        return fmt.Errorf("failed to get article: %w", err)
    }
    
    // 既に要約が存在する場合はスキップ
    if article.Summary != nil && *article.Summary != "" {
        return nil
    }
    
    // 要約生成リクエストの作成
    summaryType := "medium"
    if st, ok := payload.Options["summary_type"].(string); ok {
        summaryType = st
    }
    
    req := &SummaryRequest{
        Content:     *article.Content,
        Title:       article.Title,
        URL:         article.URL,
        Language:    article.Language,
        SummaryType: summaryType,
    }
    
    // 要約生成
    summary, err := s.aiService.GenerateSummary(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to generate summary: %w", err)
    }
    
    // 記事の更新
    article.Summary = &summary.Summary
    article.SummaryGenerationStatus = "completed"
    article.SummaryGeneratedAt = &summary.GeneratedAt
    article.SummaryModelVersion = &summary.ModelVersion
    
    if err := s.articleRepo.Update(article); err != nil {
        return fmt.Errorf("failed to update article: %w", err)
    }
    
    log.Printf("Summary generated for article %s using %s", article.ID, summary.Provider)
    return nil
}

func (s *JobService) StartWorker(ctx context.Context, workerID int) {
    log.Printf("Starting worker %d", workerID)
    
    for {
        select {
        case <-ctx.Done():
            log.Printf("Worker %d shutting down", workerID)
            return
        default:
            job, err := s.jobRepo.GetNextJob()
            if err != nil {
                time.Sleep(1 * time.Second)
                continue
            }
            
            if job == nil {
                time.Sleep(5 * time.Second)
                continue
            }
            
            s.processJobWithRetry(ctx, job)
        }
    }
}

func (s *JobService) processJobWithRetry(ctx context.Context, job *models.JobQueue) {
    job.Status = "processing"
    job.StartedAt = timePtr(time.Now())
    s.jobRepo.Update(job)
    
    err := s.ProcessJob(ctx, job)
    
    if err != nil {
        log.Printf("Job %s failed: %v", job.ID, err)
        
        job.RetryCount++
        job.ErrorMessage = stringPtr(err.Error())
        
        if job.RetryCount >= job.MaxRetries {
            job.Status = "failed"
        } else {
            job.Status = "pending"
        }
    } else {
        job.Status = "completed"
        job.CompletedAt = timePtr(time.Now())
    }
    
    s.jobRepo.Update(job)
}

func timePtr(t time.Time) *time.Time {
    return &t
}

func stringPtr(s string) *string {
    return &s
}
```

## 受入条件

### 必須条件
- [ ] Groq APIとの連携が正常に動作する
- [ ] Anthropic Claude APIとの連携が正常に動作する
- [ ] フォールバック機能が正常に動作する
- [ ] 非同期ジョブ処理が正常に動作する
- [ ] 要約生成の品質が許容レベル
- [ ] エラーハンドリングが適切に動作する

### 品質条件
- [ ] API呼び出しのレート制限を遵守
- [ ] レスポンス時間が30秒以内
- [ ] 同時処理数の制限が機能している
- [ ] ログが適切に出力される
- [ ] リトライ機能が正常に動作する

## 推定時間
**40時間** (7-10日)

## 依存関係
- P1-201: バックエンド基盤構築
- 記事保存機能の基本実装

## 完了後の次ステップ
1. P1-302: 記事管理API拡張
2. 要約品質の改善
3. カテゴリ自動分類機能

## 備考
- APIキーは環境変数で管理
- レート制限を適切に設定
- エラー時の適切なフォールバック
- 要約品質の継続的な改善
- コスト最適化を意識した実装