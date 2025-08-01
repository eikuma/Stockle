package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/eikuma/stockle/backend/internal/config"
	"github.com/eikuma/stockle/backend/internal/models"
	"github.com/eikuma/stockle/backend/internal/repositories"
)

type AuthService struct {
	userRepo  repositories.UserRepository
	jwtConfig *config.JWTConfig
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

func NewAuthService(userRepo repositories.UserRepository, jwtConfig *config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
	}
}

func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
	// メールアドレスの重複チェック
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	// パスワードハッシュ化
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: stringPtr(string(hashedPassword)),
		Name:         req.DisplayName,
		DisplayName:  req.DisplayName,
		AuthProvider: "email",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (*TokenPair, *models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if user.PasswordHash == nil {
		return nil, nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// 最終ログイン時刻の更新
	now := time.Now()
	user.LastLoginAt = &now
	s.userRepo.Update(user)

	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokens, user, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (*TokenPair, error) {
	// リフレッシュトークンの検証
	session, err := s.userRepo.GetSessionByToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("refresh token expired")
	}

	user, err := s.userRepo.GetByID(session.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// 新しいトークンペアの生成
	tokens, err := s.generateTokenPair(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// 古いセッションを削除
	s.userRepo.DeleteSession(session.ID)

	return tokens, nil
}

func (s *AuthService) generateTokenPair(user *models.User) (*TokenPair, error) {
	// アクセストークンの生成
	accessClaims := &Claims{
		UserID: fmt.Sprintf("%d", user.ID),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtConfig.AccessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.jwtConfig.Issuer,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtConfig.AccessSecret))
	if err != nil {
		return nil, err
	}

	// リフレッシュトークンの生成
	refreshToken, err := s.generateRandomToken()
	if err != nil {
		return nil, err
	}

	// セッションの保存
	session := &models.UserSession{
		UserID:     user.ID,
		TokenHash:  s.hashToken(refreshToken),
		ExpiresAt:  time.Now().Add(s.jwtConfig.RefreshExpiry),
		LastUsedAt: time.Now(),
	}

	if err := s.userRepo.CreateSession(session); err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.jwtConfig.AccessExpiry.Seconds()),
	}, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtConfig.AccessSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) Logout(refreshToken string) error {
	session, err := s.userRepo.GetSessionByToken(refreshToken)
	if err != nil {
		return errors.New("invalid refresh token")
	}

	return s.userRepo.DeleteSession(session.ID)
}

func (s *AuthService) generateRandomToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *AuthService) hashToken(token string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
	return string(hash)
}

func stringPtr(s string) *string {
	return &s
}

// リクエスト/レスポンス構造体
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"display_name" validate:"required,min=1,max=100"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}