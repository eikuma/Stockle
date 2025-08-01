package models

import (
	"time"
)

type Article struct {
	ID                      string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID                  string     `json:"userId" gorm:"not null;type:varchar(36);index"`
	CategoryID              *string    `json:"categoryId,omitempty" gorm:"type:varchar(36);index"`
	URL                     string     `json:"url" gorm:"not null;type:text"`
	Title                   string     `json:"title" gorm:"not null;type:varchar(500)"`
	Content                 *string    `json:"content,omitempty" gorm:"type:longtext"`
	Summary                 *string    `json:"summary,omitempty" gorm:"type:text"`
	SummaryShort            *string    `json:"summaryShort,omitempty" gorm:"type:text"`
	SummaryLong             *string    `json:"summaryLong,omitempty" gorm:"type:longtext"`
	ThumbnailURL            *string    `json:"thumbnailUrl,omitempty" gorm:"type:text"`
	Author                  *string    `json:"author,omitempty" gorm:"type:varchar(255)"`
	SiteName                *string    `json:"siteName,omitempty" gorm:"type:varchar(255)"`
	PublishedAt             *time.Time `json:"publishedAt,omitempty"`
	SavedAt                 time.Time  `json:"savedAt" gorm:"autoCreateTime"`
	LastAccessedAt          *time.Time `json:"lastAccessedAt,omitempty"`
	Status                  string     `json:"status" gorm:"type:varchar(20);default:'unread'"`
	IsFavorite              bool       `json:"isFavorite" gorm:"default:false"`
	ReadingProgress         float64    `json:"readingProgress" gorm:"default:0"`
	ReadingTimeSeconds      int        `json:"readingTimeSeconds" gorm:"default:0"`
	WordCount               *int       `json:"wordCount,omitempty"`
	Language                string     `json:"language" gorm:"type:varchar(10);default:'ja'"`
	SummaryGenerationStatus string     `json:"summaryGenerationStatus" gorm:"type:varchar(20);default:'pending'"`
	SummaryGeneratedAt      *time.Time `json:"summaryGeneratedAt,omitempty"`
	SummaryModelVersion     *string    `json:"summaryModelVersion,omitempty" gorm:"type:varchar(100)"`
	CreatedAt               time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt               time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`

	// Associations
	Category *Category       `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Tags     []Tag           `json:"tags" gorm:"many2many:article_tags;"`
	User     *User           `json:"-" gorm:"foreignKey:UserID"`
}

// ArticleStatus represents possible article statuses
const (
	ArticleStatusUnread   = "unread"
	ArticleStatusRead     = "read"
	ArticleStatusArchived = "archived"
)

// SummaryGenerationStatus represents possible summary generation statuses
const (
	SummaryStatusPending   = "pending"
	SummaryStatusProcessing = "processing"
	SummaryStatusCompleted = "completed"
	SummaryStatusFailed    = "failed"
)