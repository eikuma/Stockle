package repositories

import (
	"strings"

	"github.com/eikuma/stockle/backend/internal/models"
	"gorm.io/gorm"
)

// ArticleFilters represents filtering options for articles
type ArticleFilters struct {
	Status     string
	CategoryID string
	Search     string
	Favorite   *bool
	Page       int
	Limit      int
}

// ArticleListResult represents paginated article results
type ArticleListResult struct {
	Articles []*models.Article
	Total    int64
	Page     int
	Limit    int
}

type ArticleRepository interface {
	Create(article *models.Article) error
	CreateWithTags(article *models.Article, tagIDs []string) error
	Update(article *models.Article) error
	UpdateStatus(id, userID, status string) error
	UpdateFavorite(id, userID string, isFavorite bool) error
	UpdateReadingProgress(id, userID string, progress float64) error
	GetByID(id string) (*models.Article, error)
	GetByIDWithAssociations(id string) (*models.Article, error)
	GetByUserID(userID string) ([]*models.Article, error)
	GetByUserIDWithFilters(userID string, filters ArticleFilters) (*ArticleListResult, error)
	GetByURL(userID, url string) (*models.Article, error)
	Delete(id, userID string) error
	SearchByTitle(userID, query string) ([]*models.Article, error)
	GetFavorites(userID string, page, limit int) (*ArticleListResult, error)
	GetRecentlyRead(userID string, limit int) ([]*models.Article, error)
	MarkAsAccessed(id string) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{
		db: db,
	}
}

func (r *articleRepository) Create(article *models.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) CreateWithTags(article *models.Article, tagIDs []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create article
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// Associate tags if any
		if len(tagIDs) > 0 {
			var tags []models.Tag
			if err := tx.Where("id IN ? AND user_id = ?", tagIDs, article.UserID).Find(&tags).Error; err != nil {
				return err
			}
			
			if err := tx.Model(article).Association("Tags").Append(tags); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *articleRepository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) UpdateStatus(id, userID, status string) error {
	return r.db.Model(&models.Article{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("status", status).Error
}

func (r *articleRepository) UpdateFavorite(id, userID string, isFavorite bool) error {
	return r.db.Model(&models.Article{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_favorite", isFavorite).Error
}

func (r *articleRepository) UpdateReadingProgress(id, userID string, progress float64) error {
	return r.db.Model(&models.Article{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"reading_progress": progress,
			"last_accessed_at": gorm.Expr("NOW()"),
		}).Error
}

func (r *articleRepository) GetByID(id string) (*models.Article, error) {
	var article models.Article
	err := r.db.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetByIDWithAssociations(id string) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("Category").
		Preload("Tags").
		Where("id = ?", id).
		First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetByUserID(userID string) ([]*models.Article, error) {
	var articles []*models.Article
	err := r.db.Preload("Category").
		Preload("Tags").
		Where("user_id = ?", userID).
		Order("saved_at DESC").
		Find(&articles).Error
	return articles, err
}

func (r *articleRepository) GetByUserIDWithFilters(userID string, filters ArticleFilters) (*ArticleListResult, error) {
	query := r.db.Model(&models.Article{}).
		Preload("Category").
		Preload("Tags").
		Where("user_id = ?", userID)

	// Apply filters
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	if filters.CategoryID != "" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}

	if filters.Favorite != nil {
		query = query.Where("is_favorite = ?", *filters.Favorite)
	}

	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where(
			"LOWER(title) LIKE ? OR LOWER(content) LIKE ? OR LOWER(summary) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Apply pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.Limit < 1 {
		filters.Limit = 20
	}

	offset := (filters.Page - 1) * filters.Limit
	query = query.Offset(offset).Limit(filters.Limit)

	// Get articles
	var articles []*models.Article
	err := query.Order("saved_at DESC").Find(&articles).Error
	if err != nil {
		return nil, err
	}

	return &ArticleListResult{
		Articles: articles,
		Total:    total,
		Page:     filters.Page,
		Limit:    filters.Limit,
	}, nil
}

func (r *articleRepository) GetByURL(userID, url string) (*models.Article, error) {
	var article models.Article
	err := r.db.Where("user_id = ? AND url = ?", userID, url).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) Delete(id, userID string) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Article{}).Error
}

func (r *articleRepository) SearchByTitle(userID, query string) ([]*models.Article, error) {
	var articles []*models.Article
	searchTerm := "%" + strings.ToLower(query) + "%"
	
	err := r.db.Preload("Category").
		Preload("Tags").
		Where("user_id = ? AND LOWER(title) LIKE ?", userID, searchTerm).
		Order("saved_at DESC").
		Limit(50).
		Find(&articles).Error
	
	return articles, err
}

func (r *articleRepository) GetFavorites(userID string, page, limit int) (*ArticleListResult, error) {
	filters := ArticleFilters{
		Favorite: boolPtr(true),
		Page:     page,
		Limit:    limit,
	}
	return r.GetByUserIDWithFilters(userID, filters)
}

func (r *articleRepository) GetRecentlyRead(userID string, limit int) ([]*models.Article, error) {
	var articles []*models.Article
	err := r.db.Preload("Category").
		Preload("Tags").
		Where("user_id = ? AND status = ? AND last_accessed_at IS NOT NULL", userID, models.ArticleStatusRead).
		Order("last_accessed_at DESC").
		Limit(limit).
		Find(&articles).Error
	
	return articles, err
}

func (r *articleRepository) MarkAsAccessed(id string) error {
	return r.db.Model(&models.Article{}).
		Where("id = ?", id).
		Update("last_accessed_at", gorm.Expr("NOW()")).Error
}

// Helper function
func boolPtr(b bool) *bool {
	return &b
}