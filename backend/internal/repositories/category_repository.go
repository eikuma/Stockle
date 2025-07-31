package repositories

import (
	"github.com/eikuma/stockle/backend/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *models.Category) error
	Update(category *models.Category) error
	GetByID(id string) (*models.Category, error)
	GetByUserID(userID string) ([]*models.Category, error)
	GetByUserIDWithCounts(userID string) ([]*models.Category, error)
	Delete(id, userID string) error
	GetDefault(userID string) (*models.Category, error)
	CreateDefault(userID string) (*models.Category, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *models.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) GetByID(id string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetByUserID(userID string) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Where("user_id = ?", userID).
		Order("display_order ASC, created_at ASC").
		Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetByUserIDWithCounts(userID string) ([]*models.Category, error) {
	var categories []*models.Category
	err := r.db.Raw(`
		SELECT c.*, COALESCE(ac.article_count, 0) as article_count
		FROM categories c
		LEFT JOIN (
			SELECT category_id, COUNT(*) as article_count
			FROM articles
			WHERE user_id = ?
			GROUP BY category_id
		) ac ON c.id = ac.category_id
		WHERE c.user_id = ?
		ORDER BY c.display_order ASC, c.created_at ASC
	`, userID, userID).Scan(&categories).Error
	
	return categories, err
}

func (r *categoryRepository) Delete(id, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get default category
		var defaultCategory models.Category
		err := tx.Where("user_id = ? AND is_default = ?", userID, true).First(&defaultCategory).Error
		if err != nil {
			return err
		}

		// Move articles to default category
		err = tx.Model(&models.Article{}).
			Where("category_id = ? AND user_id = ?", id, userID).
			Update("category_id", defaultCategory.ID).Error
		if err != nil {
			return err
		}

		// Delete category
		return tx.Where("id = ? AND user_id = ? AND is_default = ?", id, userID, false).
			Delete(&models.Category{}).Error
	})
}

func (r *categoryRepository) GetDefault(userID string) (*models.Category, error) {
	var category models.Category
	err := r.db.Where("user_id = ? AND is_default = ?", userID, true).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) CreateDefault(userID string) (*models.Category, error) {
	category := &models.Category{
		UserID:       userID,
		Name:         "未分類",
		Color:        "#6B7280",
		DisplayOrder: 0,
		IsDefault:    true,
	}
	
	err := r.db.Create(category).Error
	if err != nil {
		return nil, err
	}
	
	return category, nil
}