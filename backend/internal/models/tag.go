package models

import (
	"time"
)

type Tag struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID     string    `json:"userId" gorm:"not null;type:varchar(36);index"`
	Name       string    `json:"name" gorm:"not null;type:varchar(50);uniqueIndex:idx_user_tag_name"`
	UsageCount int       `json:"usageCount" gorm:"default:0"`
	CreatedAt  time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Associations
	Articles []Article `json:"-" gorm:"many2many:article_tags;"`
	User     *User     `json:"-" gorm:"foreignKey:UserID"`
}

// ArticleTag represents the join table for article-tag relationship
type ArticleTag struct {
	ArticleID string    `gorm:"primaryKey;type:varchar(36)"`
	TagID     string    `gorm:"primaryKey;type:varchar(36)"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}