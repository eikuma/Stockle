package models

import (
	"time"
)

type Article struct {
	ID                      string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	UserID                  string     `json:"user_id" gorm:"not null;type:varchar(36);index"`
	Title                   string     `json:"title" gorm:"not null;type:varchar(500)"`
	URL                     string     `json:"url" gorm:"not null;type:text"`
	Content                 *string    `json:"content,omitempty" gorm:"type:longtext"`
	Summary                 *string    `json:"summary,omitempty" gorm:"type:text"`
	Language                string     `json:"language" gorm:"type:varchar(10);default:'ja'"`
	SummaryGenerationStatus string     `json:"summary_generation_status" gorm:"type:varchar(20);default:'pending'"`
	SummaryGeneratedAt      *time.Time `json:"summary_generated_at,omitempty"`
	SummaryModelVersion     *string    `json:"summary_model_version,omitempty" gorm:"type:varchar(100)"`
	CreatedAt               time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt               time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}

// SummaryGenerationStatus represents possible summary generation statuses
const (
	SummaryStatusPending   = "pending"
	SummaryStatusProcessing = "processing"
	SummaryStatusCompleted = "completed"
	SummaryStatusFailed    = "failed"
)