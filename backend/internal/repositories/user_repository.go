package repositories

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/eikuma/stockle/backend/internal/models"
)

type UserRepository interface {
	Create(user *models.User) error
	Update(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	Delete(id uint) error
	
	// セッション管理
	CreateSession(session *models.UserSession) error
	GetSessionByToken(token string) (*models.UserSession, error)
	GetSessionsByUserID(userID uint) ([]models.UserSession, error)
	UpdateSession(session *models.UserSession) error
	DeleteSession(sessionID uint) error
	DeleteExpiredSessions() error
	DeleteUserSessions(userID uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Preferences").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Preferences").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByGoogleID(googleID string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Preferences").Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// セッション管理メソッド

func (r *userRepository) CreateSession(session *models.UserSession) error {
	return r.db.Create(session).Error
}

func (r *userRepository) GetSessionByToken(token string) (*models.UserSession, error) {
	var sessions []models.UserSession
	err := r.db.Where("is_active = ?", true).Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	// トークンハッシュと比較
	for _, session := range sessions {
		if bcrypt.CompareHashAndPassword([]byte(session.TokenHash), []byte(token)) == nil {
			// セッションの有効期限チェック
			if time.Now().After(session.ExpiresAt) {
				// 期限切れセッションを無効化
				session.IsActive = false
				r.db.Save(&session)
				return nil, errors.New("session expired")
			}
			return &session, nil
		}
	}

	return nil, errors.New("session not found")
}

func (r *userRepository) GetSessionsByUserID(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).Find(&sessions).Error
	return sessions, err
}

func (r *userRepository) UpdateSession(session *models.UserSession) error {
	return r.db.Save(session).Error
}

func (r *userRepository) DeleteSession(sessionID uint) error {
	return r.db.Delete(&models.UserSession{}, sessionID).Error
}

func (r *userRepository) DeleteExpiredSessions() error {
	now := time.Now()
	return r.db.Where("expires_at < ?", now).Delete(&models.UserSession{}).Error
}

func (r *userRepository) DeleteUserSessions(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.UserSession{}).Error
}