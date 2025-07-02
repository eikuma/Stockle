package repositories

import (
	"github.com/eikuma/stockle/backend/internal/models"
	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(article *models.Article) error
	Update(article *models.Article) error
	GetByID(id string) (*models.Article, error)
	GetByUserID(userID string) ([]*models.Article, error)
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

func (r *articleRepository) Update(article *models.Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) GetByID(id string) (*models.Article, error) {
	var article models.Article
	err := r.db.Where("id = ?", id).First(&article).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetByUserID(userID string) ([]*models.Article, error) {
	var articles []*models.Article
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&articles).Error
	return articles, err
}