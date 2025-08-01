package repositories

import (
	"strings"

	"github.com/eikuma/stockle/backend/internal/models"
	"gorm.io/gorm"
)

type TagRepository interface {
	Create(tag *models.Tag) error
	GetOrCreate(userID, name string) (*models.Tag, error)
	GetOrCreateMultiple(userID string, names []string) ([]*models.Tag, error)
	GetByID(id string) (*models.Tag, error)
	GetByUserID(userID string) ([]*models.Tag, error)
	GetPopularTags(userID string, limit int) ([]*models.Tag, error)
	SearchByName(userID, query string) ([]*models.Tag, error)
	UpdateUsageCount(tagID string) error
	Delete(id, userID string) error
	GetUnusedTags(userID string) ([]*models.Tag, error)
}

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepository{
		db: db,
	}
}

func (r *tagRepository) Create(tag *models.Tag) error {
	return r.db.Create(tag).Error
}

func (r *tagRepository) GetOrCreate(userID, name string) (*models.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, gorm.ErrRecordNotFound
	}

	var tag models.Tag
	err := r.db.Where("user_id = ? AND name = ?", userID, name).First(&tag).Error
	
	if err == gorm.ErrRecordNotFound {
		// Create new tag
		tag = models.Tag{
			UserID:     userID,
			Name:       name,
			UsageCount: 0,
		}
		
		if err := r.db.Create(&tag).Error; err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *tagRepository) GetOrCreateMultiple(userID string, names []string) ([]*models.Tag, error) {
	var tags []*models.Tag
	
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		
		tag, err := r.GetOrCreate(userID, name)
		if err != nil {
			return nil, err
		}
		
		tags = append(tags, tag)
	}
	
	return tags, nil
}

func (r *tagRepository) GetByID(id string) (*models.Tag, error) {
	var tag models.Tag
	err := r.db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepository) GetByUserID(userID string) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.Where("user_id = ?", userID).
		Order("usage_count DESC, name ASC").
		Find(&tags).Error
	return tags, err
}

func (r *tagRepository) GetPopularTags(userID string, limit int) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.Where("user_id = ? AND usage_count > 0", userID).
		Order("usage_count DESC, name ASC").
		Limit(limit).
		Find(&tags).Error
	return tags, err
}

func (r *tagRepository) SearchByName(userID, query string) ([]*models.Tag, error) {
	var tags []*models.Tag
	searchTerm := "%" + strings.ToLower(strings.TrimSpace(query)) + "%"
	
	err := r.db.Where("user_id = ? AND LOWER(name) LIKE ?", userID, searchTerm).
		Order("usage_count DESC, name ASC").
		Limit(20).
		Find(&tags).Error
	
	return tags, err
}

func (r *tagRepository) UpdateUsageCount(tagID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Get current count from article_tags
		var count int64
		err := tx.Table("article_tags").Where("tag_id = ?", tagID).Count(&count).Error
		if err != nil {
			return err
		}

		// Update usage count
		return tx.Model(&models.Tag{}).
			Where("id = ?", tagID).
			Update("usage_count", count).Error
	})
}

func (r *tagRepository) Delete(id, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete article-tag associations
		err := tx.Exec("DELETE FROM article_tags WHERE tag_id = ?", id).Error
		if err != nil {
			return err
		}

		// Delete tag
		return tx.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Tag{}).Error
	})
}

func (r *tagRepository) GetUnusedTags(userID string) ([]*models.Tag, error) {
	var tags []*models.Tag
	err := r.db.Where("user_id = ? AND usage_count = 0", userID).
		Order("created_at DESC").
		Find(&tags).Error
	return tags, err
}