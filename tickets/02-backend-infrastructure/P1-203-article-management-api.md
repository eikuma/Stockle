# P1-203: 記事管理API実装

## 概要
記事の保存・取得・更新・削除・検索のAPIエンドポイント実装

## 担当者
**Member 2 (Backend Infrastructure Developer)**

## 優先度
**最高** - フロントエンドとの連携に必要

## 前提条件
- P1-201: バックエンド基盤構築が完了済み
- P1-202: 認証システム実装が完了済み
- Member 3のAI統合と並行作業可能

## 作業内容

### 1. 記事モデルとリポジトリの実装
- [ ] 記事データモデルの詳細実装
- [ ] カテゴリモデルの実装
- [ ] タグモデルの実装
- [ ] 記事リポジトリの実装
- [ ] 検索・フィルタリング機能

### 2. 記事保存API
- [ ] URL解析・メタデータ取得
- [ ] 記事コンテンツのスクレイピング
- [ ] 重複チェック機能
- [ ] バリデーション機能
- [ ] 非同期処理対応

### 3. 記事取得・一覧API
- [ ] ページネーション対応
- [ ] ソート機能
- [ ] フィルタリング機能
- [ ] 検索機能
- [ ] カテゴリ別取得

### 4. 記事更新・削除API
- [ ] 記事情報の更新
- [ ] ステータス変更
- [ ] カテゴリ変更
- [ ] タグ管理
- [ ] 論理削除

### 5. カテゴリ管理API
- [ ] カテゴリCRUD操作
- [ ] カテゴリ並び順管理
- [ ] 記事数の自動更新
- [ ] デフォルトカテゴリ設定

### 6. Web スクレイピング機能
- [ ] メタデータ抽出
- [ ] OGP情報取得
- [ ] 本文抽出
- [ ] 画像URL取得
- [ ] エラーハンドリング

## 実装詳細

### internal/models/article.go (拡張)
```go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type Article struct {
    ID                       string     `json:"id" gorm:"type:char(36);primaryKey"`
    UserID                   string     `json:"user_id" gorm:"not null;index"`
    CategoryID               *string    `json:"category_id" gorm:"type:char(36)"`
    URL                      string     `json:"url" gorm:"type:text;not null"`
    URLHash                  string     `json:"-" gorm:"type:varchar(64);not null;index"`
    Title                    string     `json:"title" gorm:"type:text;not null"`
    Content                  *string    `json:"content" gorm:"type:mediumtext"`
    ContentHash              *string    `json:"-" gorm:"type:varchar(64)"`
    Summary                  *string    `json:"summary" gorm:"type:text"`
    SummaryShort             *string    `json:"summary_short" gorm:"type:text"`
    SummaryLong              *string    `json:"summary_long" gorm:"type:text"`
    ThumbnailURL             *string    `json:"thumbnail_url" gorm:"type:text"`
    Author                   *string    `json:"author"`
    SiteName                 *string    `json:"site_name"`
    PublishedAt              *time.Time `json:"published_at"`
    SavedAt                  time.Time  `json:"saved_at" gorm:"not null;default:CURRENT_TIMESTAMP"`
    LastAccessedAt           *time.Time `json:"last_accessed_at"`
    Status                   string     `json:"status" gorm:"default:'unread';index"`
    IsFavorite               bool       `json:"is_favorite" gorm:"default:false"`
    ReadingProgress          float64    `json:"reading_progress" gorm:"default:0.00"`
    ReadingTimeSeconds       int        `json:"reading_time_seconds" gorm:"default:0"`
    WordCount                *int       `json:"word_count"`
    Language                 string     `json:"language" gorm:"default:'ja'"`
    CategoryConfidenceScore  *float64   `json:"category_confidence_score"`
    SummaryGenerationStatus  string     `json:"summary_generation_status" gorm:"default:'pending'"`
    SummaryGeneratedAt       *time.Time `json:"summary_generated_at"`
    SummaryModelVersion      *string    `json:"summary_model_version"`
    Metadata                 string     `json:"metadata" gorm:"type:json"`
    CreatedAt                time.Time  `json:"created_at"`
    UpdatedAt                time.Time  `json:"updated_at"`
    DeletedAt                gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Relations
    User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
    Tags     []Tag     `json:"tags,omitempty" gorm:"many2many:article_tags;"`
}

func (a *Article) BeforeCreate(tx *gorm.DB) error {
    if a.ID == "" {
        a.ID = uuid.New().String()
    }
    return nil
}

type Category struct {
    ID           string    `json:"id" gorm:"type:char(36);primaryKey"`
    UserID       string    `json:"user_id" gorm:"not null;index"`
    Name         string    `json:"name" gorm:"not null"`
    Color        string    `json:"color" gorm:"default:'#6B7280'"`
    DisplayOrder int       `json:"display_order" gorm:"default:0"`
    IsDefault    bool      `json:"is_default" gorm:"default:false"`
    ArticleCount int       `json:"article_count" gorm:"default:0"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`
    DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Relations
    User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Articles []Article `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) error {
    if c.ID == "" {
        c.ID = uuid.New().String()
    }
    return nil
}

type Tag struct {
    ID         string    `json:"id" gorm:"type:char(36);primaryKey"`
    UserID     string    `json:"user_id" gorm:"not null;index"`
    Name       string    `json:"name" gorm:"not null"`
    UsageCount int       `json:"usage_count" gorm:"default:0"`
    CreatedAt  time.Time `json:"created_at"`
    UpdatedAt  time.Time `json:"updated_at"`
    
    // Relations
    User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
    Articles []Article `json:"articles,omitempty" gorm:"many2many:article_tags;"`
}

func (t *Tag) BeforeCreate(tx *gorm.DB) error {
    if t.ID == "" {
        t.ID = uuid.New().String()
    }
    return nil
}

type ArticleTag struct {
    ArticleID string    `json:"article_id" gorm:"primaryKey"`
    TagID     string    `json:"tag_id" gorm:"primaryKey"`
    CreatedAt time.Time `json:"created_at"`
    
    // Relations
    Article Article `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
    Tag     Tag     `json:"tag,omitempty" gorm:"foreignKey:TagID"`
}
```

### internal/repositories/article_repository.go
```go
package repositories

import (
    "fmt"
    "gorm.io/gorm"
    
    "github.com/eikuma/stockle/backend/internal/models"
)

type ArticleRepository interface {
    Create(article *models.Article) error
    GetByID(id string) (*models.Article, error)
    GetByUserID(userID string, filters ArticleFilters) ([]*models.Article, int64, error)
    Update(article *models.Article) error
    Delete(id string) error
    GetByURLHash(userID, urlHash string) (*models.Article, error)
    Search(userID, query string, filters ArticleFilters) ([]*models.Article, int64, error)
}

type ArticleFilters struct {
    Status     string
    CategoryID string
    IsFavorite *bool
    Search     string
    Offset     int
    Limit      int
    SortBy     string
    SortOrder  string
}

type articleRepository struct {
    db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
    return &articleRepository{db: db}
}

func (r *articleRepository) Create(article *models.Article) error {
    return r.db.Create(article).Error
}

func (r *articleRepository) GetByID(id string) (*models.Article, error) {
    var article models.Article
    err := r.db.Preload("Category").Preload("Tags").First(&article, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &article, nil
}

func (r *articleRepository) GetByUserID(userID string, filters ArticleFilters) ([]*models.Article, int64, error) {
    query := r.db.Model(&models.Article{}).Where("user_id = ?", userID)
    
    // フィルタリング
    if filters.Status != "" {
        query = query.Where("status = ?", filters.Status)
    }
    if filters.CategoryID != "" {
        query = query.Where("category_id = ?", filters.CategoryID)
    }
    if filters.IsFavorite != nil {
        query = query.Where("is_favorite = ?", *filters.IsFavorite)
    }
    if filters.Search != "" {
        searchTerm := "%" + filters.Search + "%"
        query = query.Where("title LIKE ? OR summary LIKE ?", searchTerm, searchTerm)
    }
    
    // 総数を取得
    var total int64
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // ソート
    if filters.SortBy != "" {
        order := filters.SortBy
        if filters.SortOrder == "desc" {
            order += " DESC"
        }
        query = query.Order(order)
    } else {
        query = query.Order("saved_at DESC")
    }
    
    // ページネーション
    if filters.Limit > 0 {
        query = query.Limit(filters.Limit)
    }
    if filters.Offset > 0 {
        query = query.Offset(filters.Offset)
    }
    
    var articles []*models.Article
    err := query.Preload("Category").Preload("Tags").Find(&articles).Error
    return articles, total, err
}

func (r *articleRepository) Update(article *models.Article) error {
    return r.db.Save(article).Error
}

func (r *articleRepository) Delete(id string) error {
    return r.db.Delete(&models.Article{}, "id = ?", id).Error
}

func (r *articleRepository) GetByURLHash(userID, urlHash string) (*models.Article, error) {
    var article models.Article
    err := r.db.Where("user_id = ? AND url_hash = ?", userID, urlHash).First(&article).Error
    if err != nil {
        return nil, err
    }
    return &article, nil
}

func (r *articleRepository) Search(userID, query string, filters ArticleFilters) ([]*models.Article, int64, error) {
    // 全文検索の実装
    searchQuery := r.db.Model(&models.Article{}).Where("user_id = ?", userID)
    
    if query != "" {
        searchQuery = searchQuery.Where(
            "MATCH(title, content, summary) AGAINST(? IN NATURAL LANGUAGE MODE)", 
            query,
        )
    }
    
    // その他のフィルタリング
    if filters.Status != "" {
        searchQuery = searchQuery.Where("status = ?", filters.Status)
    }
    if filters.CategoryID != "" {
        searchQuery = searchQuery.Where("category_id = ?", filters.CategoryID)
    }
    
    var total int64
    if err := searchQuery.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    var articles []*models.Article
    err := searchQuery.
        Offset(filters.Offset).
        Limit(filters.Limit).
        Order("saved_at DESC").
        Preload("Category").
        Preload("Tags").
        Find(&articles).Error
    
    return articles, total, err
}
```

### internal/services/scraper_service.go
```go
package services

import (
    "fmt"
    "net/http"
    "net/url"
    "regexp"
    "strings"
    "time"
    
    "github.com/PuerkitoBio/goquery"
)

type ScraperService struct {
    client *http.Client
}

type ScrapedData struct {
    Title        string
    Content      string
    Description  string
    Author       string
    SiteName     string
    ThumbnailURL string
    PublishedAt  *time.Time
    Language     string
    WordCount    int
}

func NewScraperService() *ScraperService {
    return &ScraperService{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (s *ScraperService) ScrapeURL(rawURL string) (*ScrapedData, error) {
    // URL検証
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        return nil, fmt.Errorf("invalid URL: %w", err)
    }
    
    if parsedURL.Scheme == "" {
        parsedURL.Scheme = "https"
    }
    
    // HTTPリクエスト
    req, err := http.NewRequest("GET", parsedURL.String(), nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Stockle/1.0)")
    
    resp, err := s.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to fetch URL: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
    }
    
    // HTML解析
    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to parse HTML: %w", err)
    }
    
    data := &ScrapedData{}
    
    // タイトル取得
    data.Title = s.extractTitle(doc)
    
    // 説明文取得
    data.Description = s.extractDescription(doc)
    
    // 本文取得
    data.Content = s.extractContent(doc)
    
    // 著者取得
    data.Author = s.extractAuthor(doc)
    
    // サイト名取得
    data.SiteName = s.extractSiteName(doc, parsedURL)
    
    // サムネイル取得
    data.ThumbnailURL = s.extractThumbnail(doc, parsedURL)
    
    // 言語取得
    data.Language = s.extractLanguage(doc)
    
    // 単語数計算
    data.WordCount = s.countWords(data.Content)
    
    return data, nil
}

func (s *ScraperService) extractTitle(doc *goquery.Document) string {
    // OGP title
    if title, exists := doc.Find("meta[property='og:title']").Attr("content"); exists && title != "" {
        return strings.TrimSpace(title)
    }
    
    // Twitter title
    if title, exists := doc.Find("meta[name='twitter:title']").Attr("content"); exists && title != "" {
        return strings.TrimSpace(title)
    }
    
    // HTML title
    title := doc.Find("title").Text()
    return strings.TrimSpace(title)
}

func (s *ScraperService) extractDescription(doc *goquery.Document) string {
    // OGP description
    if desc, exists := doc.Find("meta[property='og:description']").Attr("content"); exists && desc != "" {
        return strings.TrimSpace(desc)
    }
    
    // Meta description
    if desc, exists := doc.Find("meta[name='description']").Attr("content"); exists && desc != "" {
        return strings.TrimSpace(desc)
    }
    
    return ""
}

func (s *ScraperService) extractContent(doc *goquery.Document) string {
    // 一般的な記事コンテナを試行
    selectors := []string{
        "article",
        "[role='main']",
        ".content",
        ".post-content",
        ".entry-content",
        ".article-body",
        "main",
    }
    
    for _, selector := range selectors {
        if content := doc.Find(selector).First(); content.Length() > 0 {
            text := content.Text()
            if len(text) > 200 { // 十分な長さがある場合
                return s.cleanText(text)
            }
        }
    }
    
    // フォールバック: body全体から不要要素を除去
    content := doc.Find("body").Clone()
    content.Find("script, style, nav, header, footer, aside, .sidebar, .menu").Remove()
    
    return s.cleanText(content.Text())
}

func (s *ScraperService) extractAuthor(doc *goquery.Document) string {
    // OGP author
    if author, exists := doc.Find("meta[property='article:author']").Attr("content"); exists && author != "" {
        return strings.TrimSpace(author)
    }
    
    // Meta author
    if author, exists := doc.Find("meta[name='author']").Attr("content"); exists && author != "" {
        return strings.TrimSpace(author)
    }
    
    // 構造化データ
    if author := doc.Find("[rel='author']").Text(); author != "" {
        return strings.TrimSpace(author)
    }
    
    return ""
}

func (s *ScraperService) extractSiteName(doc *goquery.Document, parsedURL *url.URL) string {
    // OGP site_name
    if siteName, exists := doc.Find("meta[property='og:site_name']").Attr("content"); exists && siteName != "" {
        return strings.TrimSpace(siteName)
    }
    
    // フォールバック: ドメイン名
    return parsedURL.Host
}

func (s *ScraperService) extractThumbnail(doc *goquery.Document, parsedURL *url.URL) string {
    // OGP image
    if image, exists := doc.Find("meta[property='og:image']").Attr("content"); exists && image != "" {
        return s.resolveURL(image, parsedURL)
    }
    
    // Twitter image
    if image, exists := doc.Find("meta[name='twitter:image']").Attr("content"); exists && image != "" {
        return s.resolveURL(image, parsedURL)
    }
    
    return ""
}

func (s *ScraperService) extractLanguage(doc *goquery.Document) string {
    if lang, exists := doc.Find("html").Attr("lang"); exists && lang != "" {
        return lang[:2] // 言語コードの最初の2文字
    }
    return "ja" // デフォルト
}

func (s *ScraperService) cleanText(text string) string {
    // 余分な空白文字を除去
    re := regexp.MustCompile(`\s+`)
    text = re.ReplaceAllString(text, " ")
    return strings.TrimSpace(text)
}

func (s *ScraperService) countWords(text string) int {
    if text == "" {
        return 0
    }
    words := strings.Fields(text)
    return len(words)
}

func (s *ScraperService) resolveURL(rawURL string, baseURL *url.URL) string {
    parsedURL, err := url.Parse(rawURL)
    if err != nil {
        return rawURL
    }
    
    return baseURL.ResolveReference(parsedURL).String()
}
```

### internal/controllers/article_controller.go
```go
package controllers

import (
    "crypto/sha256"
    "encoding/hex"
    "net/http"
    "strconv"
    
    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
    
    "github.com/eikuma/stockle/backend/internal/services"
    "github.com/eikuma/stockle/backend/internal/repositories"
)

type ArticleController struct {
    articleRepo    repositories.ArticleRepository
    categoryRepo   repositories.CategoryRepository
    tagRepo        repositories.TagRepository
    scraperService *services.ScraperService
    jobService     *services.JobService
    validator      *validator.Validate
}

func NewArticleController(
    articleRepo repositories.ArticleRepository,
    categoryRepo repositories.CategoryRepository,
    tagRepo repositories.TagRepository,
    scraperService *services.ScraperService,
    jobService *services.JobService,
) *ArticleController {
    return &ArticleController{
        articleRepo:    articleRepo,
        categoryRepo:   categoryRepo,
        tagRepo:        tagRepo,
        scraperService: scraperService,
        jobService:     jobService,
        validator:      validator.New(),
    }
}

type SaveArticleRequest struct {
    URL        string   `json:"url" validate:"required,url"`
    CategoryID *string  `json:"category_id"`
    Tags       []string `json:"tags"`
}

func (ac *ArticleController) SaveArticle(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    var req SaveArticleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    if err := ac.validator.Struct(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // URL重複チェック
    urlHash := ac.generateURLHash(req.URL)
    existingArticle, _ := ac.articleRepo.GetByURLHash(userID, urlHash)
    if existingArticle != nil {
        c.JSON(http.StatusConflict, gin.H{
            "error": "Article already exists",
            "article": existingArticle,
        })
        return
    }
    
    // Webスクレイピング
    scrapedData, err := ac.scraperService.ScrapeURL(req.URL)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to scrape URL: " + err.Error()})
        return
    }
    
    // 記事作成
    article := &models.Article{
        UserID:       userID,
        CategoryID:   req.CategoryID,
        URL:          req.URL,
        URLHash:      urlHash,
        Title:        scrapedData.Title,
        Content:      &scrapedData.Content,
        ThumbnailURL: &scrapedData.ThumbnailURL,
        Author:       &scrapedData.Author,
        SiteName:     &scrapedData.SiteName,
        WordCount:    &scrapedData.WordCount,
        Language:     scrapedData.Language,
        Status:       "unread",
    }
    
    if err := ac.articleRepo.Create(article); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save article"})
        return
    }
    
    // タグの追加
    if len(req.Tags) > 0 {
        for _, tagName := range req.Tags {
            tag, err := ac.tagRepo.GetOrCreate(userID, tagName)
            if err != nil {
                continue // エラーを無視してタグ追加を続行
            }
            ac.tagRepo.AddToArticle(article.ID, tag.ID)
        }
    }
    
    // 要約生成ジョブをキューに追加
    ac.jobService.EnqueueSummaryJob(article.ID, 5)
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Article saved successfully",
        "article": article,
    })
}

func (ac *ArticleController) GetArticles(c *gin.Context) {
    userID := c.GetString("user_id")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }
    
    // クエリパラメータの解析
    filters := repositories.ArticleFilters{
        Status:     c.Query("status"),
        CategoryID: c.Query("category_id"),
        Search:     c.Query("search"),
        SortBy:     c.DefaultQuery("sort_by", "saved_at"),
        SortOrder:  c.DefaultQuery("sort_order", "desc"),
    }
    
    if page := c.Query("page"); page != "" {
        if p, err := strconv.Atoi(page); err == nil && p > 0 {
            filters.Offset = (p - 1) * 20
        }
    }
    
    if limit := c.Query("limit"); limit != "" {
        if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
            filters.Limit = l
        } else {
            filters.Limit = 20
        }
    } else {
        filters.Limit = 20
    }
    
    if favorite := c.Query("favorite"); favorite == "true" {
        isFavorite := true
        filters.IsFavorite = &isFavorite
    }
    
    articles, total, err := ac.articleRepo.GetByUserID(userID, filters)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get articles"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "articles": articles,
        "total":    total,
        "page":     (filters.Offset / filters.Limit) + 1,
        "limit":    filters.Limit,
    })
}

func (ac *ArticleController) GetArticle(c *gin.Context) {
    userID := c.GetString("user_id")
    articleID := c.Param("id")
    
    article, err := ac.articleRepo.GetByID(articleID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
        return
    }
    
    if article.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
    
    // 最終アクセス時刻を更新
    now := time.Now()
    article.LastAccessedAt = &now
    ac.articleRepo.Update(article)
    
    c.JSON(http.StatusOK, gin.H{"article": article})
}

func (ac *ArticleController) UpdateArticle(c *gin.Context) {
    userID := c.GetString("user_id")
    articleID := c.Param("id")
    
    article, err := ac.articleRepo.GetByID(articleID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
        return
    }
    
    if article.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
    
    var updates map[string]interface{}
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // 許可されたフィールドのみ更新
    allowedFields := map[string]bool{
        "status":           true,
        "is_favorite":      true,
        "category_id":      true,
        "reading_progress": true,
    }
    
    for field := range updates {
        if !allowedFields[field] {
            delete(updates, field)
        }
    }
    
    if err := ac.articleRepo.UpdateFields(articleID, updates); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update article"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Article updated successfully"})
}

func (ac *ArticleController) DeleteArticle(c *gin.Context) {
    userID := c.GetString("user_id")
    articleID := c.Param("id")
    
    article, err := ac.articleRepo.GetByID(articleID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
        return
    }
    
    if article.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
        return
    }
    
    if err := ac.articleRepo.Delete(articleID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func (ac *ArticleController) generateURLHash(url string) string {
    hash := sha256.Sum256([]byte(url))
    return hex.EncodeToString(hash[:])
}
```

## 受入条件

### 必須条件
- [ ] 記事保存APIが正常に動作する
- [ ] Webスクレイピングが正常に動作する
- [ ] 記事一覧・詳細取得APIが正常に動作する
- [ ] 検索・フィルタリング機能が正常に動作する
- [ ] 記事更新・削除APIが正常に動作する
- [ ] ページネーションが正常に動作する

### 品質条件
- [ ] API応答時間が200ms以内（95パーセンタイル）
- [ ] エラーハンドリングが適切
- [ ] バリデーションが適切に動作する
- [ ] セキュリティが適切に実装されている
- [ ] データベースクエリが最適化されている

## 推定時間
**40時間** (7-10日)

## 依存関係
- P1-201: バックエンド基盤構築
- P1-202: 認証システム実装
- Member 3のAI統合機能

## 完了後の次ステップ
1. フロントエンドとの統合テスト
2. パフォーマンスチューニング
3. エラーハンドリングの改善

## 備考
- セキュリティを最優先に実装
- パフォーマンスを意識したクエリ設計
- 適切なエラーハンドリング
- スクレイピング対象サイトのrobots.txtを尊重