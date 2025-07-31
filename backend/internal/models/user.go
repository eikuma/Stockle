package models

import (
	"time"
)

type User struct {
	BaseModel
	Email         string         `json:"email" gorm:"uniqueIndex;size:255;not null" validate:"required,email"`
	PasswordHash  *string        `json:"-" gorm:"size:255"`
	GoogleID      string         `json:"google_id,omitempty" gorm:"uniqueIndex;size:255"`
	Name          string         `json:"name" gorm:"size:255;not null" validate:"required,min=1,max=255"`
	DisplayName   string         `json:"display_name" gorm:"size:255;not null" validate:"required,min=1,max=255"`
	AuthProvider  string         `json:"auth_provider" gorm:"size:50;default:'email'" validate:"required"`
	AvatarURL     string         `json:"avatar_url,omitempty" gorm:"size:500"`
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	EmailVerified bool           `json:"email_verified" gorm:"default:false"`
	LastLoginAt   *time.Time     `json:"last_login_at,omitempty"`
	
	// Relationships
	Sessions    []UserSession    `json:"-" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Preferences *UserPreference  `json:"preferences,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

type UserSession struct {
	BaseModel
	UserID      uint      `json:"user_id" gorm:"not null;index"`
	TokenHash   string    `json:"-" gorm:"uniqueIndex;size:500;not null"`
	UserAgent   string    `json:"user_agent,omitempty" gorm:"size:500"`
	IPAddress   string    `json:"ip_address,omitempty" gorm:"size:45"`
	ExpiresAt   time.Time `json:"expires_at" gorm:"not null"`
	LastUsedAt  time.Time `json:"last_used_at" gorm:"not null"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	
	// Relationships
	User User `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

type UserPreference struct {
	UserID             uint   `json:"user_id" gorm:"primaryKey"`
	Language           string `json:"language" gorm:"size:10;default:'ja'"`
	Theme              string `json:"theme" gorm:"size:10;default:'light'"`
	NotificationsEmail bool   `json:"notifications_email" gorm:"default:true"`
	NotificationsPush  bool   `json:"notifications_push" gorm:"default:false"`
	AutoSummarize      bool   `json:"auto_summarize" gorm:"default:true"`
	SummaryLanguage    string `json:"summary_language" gorm:"size:10;default:'ja'"`
	
	TimestampModel
	
	// Relationships
	User User `json:"user,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

type UserCreateRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=128"`
	Name     string `json:"name" validate:"required,min=1,max=255"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserUpdateRequest struct {
	Name      string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	AvatarURL string `json:"avatar_url,omitempty" validate:"omitempty,url"`
}

type UserResponse struct {
	ID            uint                `json:"id"`
	Email         string              `json:"email"`
	Name          string              `json:"name"`
	DisplayName   string              `json:"display_name"`
	AuthProvider  string              `json:"auth_provider"`
	AvatarURL     string              `json:"avatar_url,omitempty"`
	IsActive      bool                `json:"is_active"`
	EmailVerified bool                `json:"email_verified"`
	LastLoginAt   *time.Time          `json:"last_login_at,omitempty"`
	Preferences   *UserPreference     `json:"preferences,omitempty"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		Name:          u.Name,
		DisplayName:   u.DisplayName,
		AuthProvider:  u.AuthProvider,
		AvatarURL:     u.AvatarURL,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		LastLoginAt:   u.LastLoginAt,
		Preferences:   u.Preferences,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}