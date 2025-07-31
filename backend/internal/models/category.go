package models

import (
	"time"
)

type Category struct {
	ID           string    `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID       string    `json:"userId" gorm:"not null;type:varchar(36);index"`
	Name         string    `json:"name" gorm:"not null;type:varchar(100)"`
	Color        string    `json:"color" gorm:"type:varchar(7);default:'#6B7280'"`
	DisplayOrder int       `json:"displayOrder" gorm:"default:0"`
	IsDefault    bool      `json:"isDefault" gorm:"default:false"`
	ArticleCount int       `json:"articleCount" gorm:"-"`
	CreatedAt    time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updatedAt" gorm:"autoUpdateTime"`

	// Associations
	Articles []Article `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
	User     *User     `json:"-" gorm:"foreignKey:UserID"`
}