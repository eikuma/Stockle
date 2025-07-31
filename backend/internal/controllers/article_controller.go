package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/eikuma/stockle/backend/internal/models"
	"github.com/eikuma/stockle/backend/internal/repositories"
	"github.com/eikuma/stockle/backend/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ArticleController struct {
	articleRepo  repositories.ArticleRepository
	categoryRepo repositories.CategoryRepository
	tagRepo      repositories.TagRepository
	scraperSvc   *services.ScraperService
}

type SaveArticleRequest struct {
	URL        string   `json:"url" binding:"required,url"`
	CategoryID *string  `json:"categoryId,omitempty"`
	Tags       []string `json:"tags,omitempty"`
}

type UpdateArticleRequest struct {
	Status          *string  `json:"status,omitempty"`
	IsFavorite      *bool    `json:"isFavorite,omitempty"`
	CategoryID      *string  `json:"categoryId,omitempty"`
	ReadingProgress *float64 `json:"readingProgress,omitempty"`
}

type ArticleResponse struct {
	Message string         `json:"message"`
	Article *models.Article `json:"article,omitempty"`
}

type ArticleListResponse struct {
	Articles []*models.Article `json:"articles"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Limit    int               `json:"limit"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func NewArticleController(
	articleRepo repositories.ArticleRepository,
	categoryRepo repositories.CategoryRepository,
	tagRepo repositories.TagRepository,
	scraperSvc *services.ScraperService,
) *ArticleController {
	return &ArticleController{
		articleRepo:  articleRepo,
		categoryRepo: categoryRepo,
		tagRepo:      tagRepo,
		scraperSvc:   scraperSvc,
	}
}

// SaveArticle saves a new article from URL
// POST /api/articles
func (c *ArticleController) SaveArticle(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	var req SaveArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON format: " + err.Error(),
		})
		return
	}

	// Check if article already exists
	existingArticle, _ := c.articleRepo.GetByURL(userID, req.URL)
	if existingArticle != nil {
		ctx.JSON(http.StatusConflict, ErrorResponse{
			Error:   "duplicate_article",
			Message: "This article has already been saved",
		})
		return
	}

	// Scrape article metadata
	metadata, err := c.scraperSvc.ExtractMetadata(req.URL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "scraping_failed",
			Message: "Failed to extract article content: " + err.Error(),
		})
		return
	}

	// Create article
	article := &models.Article{
		ID:           uuid.New().String(),
		UserID:       userID,
		CategoryID:   req.CategoryID,
		URL:          req.URL,
		Title:        metadata.Title,
		Content:      &metadata.Content,
		ThumbnailURL: &metadata.ThumbnailURL,
		Author:       &metadata.Author,
		SiteName:     &metadata.SiteName,
		PublishedAt:  metadata.PublishedAt,
		Status:       models.ArticleStatusUnread,
		IsFavorite:   false,
		Language:     metadata.Language,
	}

	// Create tags if provided
	var tagIDs []string
	if len(req.Tags) > 0 {
		tags, err := c.tagRepo.GetOrCreateMultiple(userID, req.Tags)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "tag_creation_failed",
				Message: "Failed to create tags: " + err.Error(),
			})
			return
		}
		
		for _, tag := range tags {
			tagIDs = append(tagIDs, tag.ID)
		}
	}

	// Save article with tags
	if err := c.articleRepo.CreateWithTags(article, tagIDs); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "save_failed",
			Message: "Failed to save article: " + err.Error(),
		})
		return
	}

	// Get article with associations for response
	savedArticle, err := c.articleRepo.GetByIDWithAssociations(article.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "fetch_failed",
			Message: "Article saved but failed to fetch details",
		})
		return
	}

	ctx.JSON(http.StatusCreated, ArticleResponse{
		Message: "Article saved successfully",
		Article: savedArticle,
	})
}

// GetArticles retrieves articles with filtering and pagination
// GET /api/articles
func (c *ArticleController) GetArticles(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	status := ctx.Query("status")
	categoryID := ctx.Query("category_id")
	search := ctx.Query("search")
	
	var favorite *bool
	if favoriteStr := ctx.Query("favorite"); favoriteStr != "" {
		if fav, err := strconv.ParseBool(favoriteStr); err == nil {
			favorite = &fav
		}
	}

	filters := repositories.ArticleFilters{
		Status:     status,
		CategoryID: categoryID,
		Search:     search,
		Favorite:   favorite,
		Page:       page,
		Limit:      limit,
	}

	result, err := c.articleRepo.GetByUserIDWithFilters(userID, filters)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "fetch_failed",
			Message: "Failed to fetch articles: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ArticleListResponse{
		Articles: result.Articles,
		Total:    result.Total,
		Page:     result.Page,
		Limit:    result.Limit,
	})
}

// GetArticle retrieves a single article by ID
// GET /api/articles/:id
func (c *ArticleController) GetArticle(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Article ID is required",
		})
		return
	}

	article, err := c.articleRepo.GetByIDWithAssociations(articleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Article not found",
		})
		return
	}

	// Check ownership
	if article.UserID != userID {
		ctx.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "Access denied",
		})
		return
	}

	// Mark as accessed
	c.articleRepo.MarkAsAccessed(articleID)

	ctx.JSON(http.StatusOK, map[string]*models.Article{
		"article": article,
	})
}

// UpdateArticle updates an existing article
// PATCH /api/articles/:id
func (c *ArticleController) UpdateArticle(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Article ID is required",
		})
		return
	}

	var req UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Invalid JSON format: " + err.Error(),
		})
		return
	}

	// Check article exists and ownership
	article, err := c.articleRepo.GetByID(articleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Article not found",
		})
		return
	}

	if article.UserID != userID {
		ctx.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "Access denied",
		})
		return
	}

	// Update fields
	if req.Status != nil {
		if err := c.articleRepo.UpdateStatus(articleID, userID, *req.Status); err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "update_failed",
				Message: "Failed to update status: " + err.Error(),
			})
			return
		}
	}

	if req.IsFavorite != nil {
		if err := c.articleRepo.UpdateFavorite(articleID, userID, *req.IsFavorite); err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "update_failed",
				Message: "Failed to update favorite status: " + err.Error(),
			})
			return
		}
	}

	if req.ReadingProgress != nil {
		if err := c.articleRepo.UpdateReadingProgress(articleID, userID, *req.ReadingProgress); err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "update_failed",
				Message: "Failed to update reading progress: " + err.Error(),
			})
			return
		}
	}

	if req.CategoryID != nil {
		article.CategoryID = req.CategoryID
		if err := c.articleRepo.Update(article); err != nil {
			ctx.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "update_failed",
				Message: "Failed to update category: " + err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusOK, ArticleResponse{
		Message: "Article updated successfully",
	})
}

// DeleteArticle deletes an article
// DELETE /api/articles/:id
func (c *ArticleController) DeleteArticle(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	articleID := ctx.Param("id")
	if articleID == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Article ID is required",
		})
		return
	}

	// Check article exists and ownership
	article, err := c.articleRepo.GetByID(articleID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "not_found",
			Message: "Article not found",
		})
		return
	}

	if article.UserID != userID {
		ctx.JSON(http.StatusForbidden, ErrorResponse{
			Error:   "forbidden",
			Message: "Access denied",
		})
		return
	}

	if err := c.articleRepo.Delete(articleID, userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "delete_failed",
			Message: "Failed to delete article: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, ArticleResponse{
		Message: "Article deleted successfully",
	})
}

// SearchArticles searches articles by title/content
// GET /api/articles/search
func (c *ArticleController) SearchArticles(ctx *gin.Context) {
	userID := ctx.GetString("user_id")
	if userID == "" {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "unauthorized",
			Message: "User not authenticated",
		})
		return
	}

	query := strings.TrimSpace(ctx.Query("q"))
	if query == "" {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "invalid_request",
			Message: "Search query is required",
		})
		return
	}

	articles, err := c.articleRepo.SearchByTitle(userID, query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "search_failed",
			Message: "Failed to search articles: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, map[string][]*models.Article{
		"articles": articles,
	})
}