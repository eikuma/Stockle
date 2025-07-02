package models

import (
	"time"
)

type JobQueue struct {
	ID           string     `json:"id" gorm:"primaryKey;type:varchar(36)"`
	JobType      string     `json:"job_type" gorm:"not null;type:varchar(50)"`
	Priority     int        `json:"priority" gorm:"not null;default:5"`
	Status       string     `json:"status" gorm:"not null;type:varchar(20);default:'pending'"`
	Payload      string     `json:"payload" gorm:"type:text"`
	MaxRetries   int        `json:"max_retries" gorm:"not null;default:3"`
	RetryCount   int        `json:"retry_count" gorm:"not null;default:0"`
	ErrorMessage *string    `json:"error_message,omitempty" gorm:"type:text"`
	CreatedAt    time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

// JobStatus represents possible job statuses
const (
	JobStatusPending    = "pending"
	JobStatusProcessing = "processing"
	JobStatusCompleted  = "completed"
	JobStatusFailed     = "failed"
)

// JobType represents possible job types
const (
	JobTypeSummarize = "summarize"
)

// JobPriority represents job priority levels
const (
	JobPriorityHigh   = 1
	JobPriorityMedium = 5
	JobPriorityLow    = 10
)