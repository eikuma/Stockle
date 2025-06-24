# P1-201: バックエンド基盤構築

## 概要
Go + Gin + GORM + MySQLによるバックエンド基盤の構築

## 担当者
**Member 2 (Backend Infrastructure Developer)**

## 優先度
**最高** - 他のバックエンド作業の基盤

## 前提条件
- P1-001: プロジェクト基盤セットアップが完了済み
- Go 1.21以上がインストール済み
- Docker環境が構築済み

## 作業内容

### 1. Go プロジェクトの初期化
- [ ] `backend/` ディレクトリに移動
- [ ] `go mod init` でモジュール初期化
- [ ] 基本的なフォルダ構造の作成
- [ ] `.gitignore` の設定
- [ ] `Makefile` の作成

### 2. 必要パッケージのインストール
- [ ] Web フレームワーク: Gin
- [ ] ORM: GORM + MySQL driver
- [ ] 認証: golang-jwt/jwt
- [ ] バリデーション: go-playground/validator
- [ ] 設定管理: Viper
- [ ] ロギング: Zap
- [ ] CORS: gin-cors
- [ ] セキュリティ: crypto/bcrypt
- [ ] ホットリロード: Air

### 3. プロジェクト構造の作成
- [ ] `cmd/api/` - メインエントリーポイント
- [ ] `internal/config/` - 設定管理
- [ ] `internal/database/` - データベース接続
- [ ] `internal/models/` - データモデル
- [ ] `internal/controllers/` - HTTPハンドラー
- [ ] `internal/middleware/` - ミドルウェア
- [ ] `internal/repositories/` - データアクセス層
- [ ] `internal/services/` - ビジネスロジック
- [ ] `internal/validators/` - バリデーション
- [ ] `internal/utils/` - ユーティリティ
- [ ] `pkg/` - 外部パッケージ
- [ ] `migrations/` - DBマイグレーション

### 4. 設定管理システム
- [ ] 環境変数の読み込み
- [ ] Viperによる設定管理
- [ ] 開発/本番環境の分離
- [ ] シークレット管理

### 5. データベース接続とモデル
- [ ] GORM設定
- [ ] MySQL接続設定
- [ ] コネクションプール設定
- [ ] 基本モデルの作成
- [ ] マイグレーション機能

### 6. ミドルウェアの実装
- [ ] CORS設定
- [ ] ロギングミドルウェア
- [ ] リカバリーミドルウェア
- [ ] レート制限
- [ ] セキュリティヘッダー

### 7. ヘルスチェックAPI
- [ ] `/health` エンドポイント
- [ ] データベース接続確認
- [ ] システム情報の返却

## 実装詳細

### go.mod
```go
module github.com/eikuma/stockle/backend

go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/gin-contrib/cors v1.4.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/mysql v1.5.2
    github.com/golang-jwt/jwt/v5 v5.0.0
    github.com/go-playground/validator/v10 v10.15.5
    github.com/spf13/viper v1.17.0
    go.uber.org/zap v1.26.0
    golang.org/x/crypto v0.14.0
    github.com/google/uuid v1.4.0
)
```

### internal/config/config.go
```go
package config

import (
    "log"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    JWT      JWTConfig      `mapstructure:"jwt"`
    AI       AIConfig       `mapstructure:"ai"`
    OAuth    OAuthConfig    `mapstructure:"oauth"`
}

type ServerConfig struct {
    Port string `mapstructure:"port"`
    Mode string `mapstructure:"mode"` // debug, release, test
}

type DatabaseConfig struct {
    Host     string `mapstructure:"host"`
    Port     string `mapstructure:"port"`
    User     string `mapstructure:"user"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"dbname"`
    SSLMode  string `mapstructure:"sslmode"`
}

type JWTConfig struct {
    Secret     string `mapstructure:"secret"`
    ExpireTime string `mapstructure:"expire_time"`
}

type AIConfig struct {
    GroqAPIKey      string `mapstructure:"groq_api_key"`
    AnthropicAPIKey string `mapstructure:"anthropic_api_key"`
}

type OAuthConfig struct {
    GoogleClientID     string `mapstructure:"google_client_id"`
    GoogleClientSecret string `mapstructure:"google_client_secret"`
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath(".")
    viper.AddConfigPath("./configs")
    
    // 環境変数の自動読み込み
    viper.AutomaticEnv()
    
    // デフォルト値の設定
    viper.SetDefault("server.port", "8080")
    viper.SetDefault("server.mode", "debug")
    viper.SetDefault("database.host", "localhost")
    viper.SetDefault("database.port", "3306")
    
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
            return nil, err
        }
        log.Println("Config file not found, using environment variables")
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### internal/database/database.go
```go
package database

import (
    "fmt"
    "log"
    "time"
    
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
    
    "github.com/eikuma/stockle/backend/internal/config"
    "github.com/eikuma/stockle/backend/internal/models"
)

type Database struct {
    DB *gorm.DB
}

func New(cfg *config.DatabaseConfig) (*Database, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.User,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.DBName,
    )
    
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })
    
    if err != nil {
        return nil, fmt.Errorf("failed to connect to database: %w", err)
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
    }
    
    // コネクションプールの設定
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return &Database{DB: db}, nil
}

func (d *Database) AutoMigrate() error {
    return d.DB.AutoMigrate(
        &models.User{},
        &models.Category{},
        &models.Article{},
        &models.Tag{},
        &models.ArticleTag{},
        &models.UserSession{},
        &models.UserPreference{},
    )
}

func (d *Database) Close() error {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Close()
}

func (d *Database) Ping() error {
    sqlDB, err := d.DB.DB()
    if err != nil {
        return err
    }
    return sqlDB.Ping()
}
```

### internal/models/user.go
```go
package models

import (
    "time"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID                        string     `json:"id" gorm:"type:char(36);primaryKey"`
    Email                     string     `json:"email" gorm:"uniqueIndex;not null"`
    PasswordHash              *string    `json:"-" gorm:"type:varchar(255)"`
    DisplayName               string     `json:"display_name" gorm:"not null"`
    ProfileImageURL           *string    `json:"profile_image_url" gorm:"type:text"`
    AuthProvider              string     `json:"auth_provider" gorm:"default:'email'"`
    GoogleID                  *string    `json:"google_id" gorm:"uniqueIndex"`
    EmailVerified             bool       `json:"email_verified" gorm:"default:false"`
    EmailVerificationToken    *string    `json:"-"`
    EmailVerificationExpiresAt *time.Time `json:"-"`
    PasswordResetToken        *string    `json:"-"`
    PasswordResetExpiresAt    *time.Time `json:"-"`
    LastLoginAt               *time.Time `json:"last_login_at"`
    CreatedAt                 time.Time  `json:"created_at"`
    UpdatedAt                 time.Time  `json:"updated_at"`
    DeletedAt                 gorm.DeletedAt `json:"-" gorm:"index"`
    
    // Relations
    Categories    []Category    `json:"categories,omitempty" gorm:"foreignKey:UserID"`
    Articles      []Article     `json:"articles,omitempty" gorm:"foreignKey:UserID"`
    Tags          []Tag         `json:"tags,omitempty" gorm:"foreignKey:UserID"`
    Sessions      []UserSession `json:"-" gorm:"foreignKey:UserID"`
    Preferences   *UserPreference `json:"preferences,omitempty" gorm:"foreignKey:UserID"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
    if u.ID == "" {
        u.ID = uuid.New().String()
    }
    return nil
}

type UserSession struct {
    ID          string     `json:"id" gorm:"type:char(36);primaryKey"`
    UserID      string     `json:"user_id" gorm:"not null"`
    TokenHash   string     `json:"-" gorm:"uniqueIndex;not null"`
    DeviceInfo  *string    `json:"device_info" gorm:"type:text"`
    IPAddress   *string    `json:"ip_address"`
    ExpiresAt   time.Time  `json:"expires_at"`
    LastUsedAt  time.Time  `json:"last_used_at"`
    CreatedAt   time.Time  `json:"created_at"`
    
    // Relations
    User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (us *UserSession) BeforeCreate(tx *gorm.DB) error {
    if us.ID == "" {
        us.ID = uuid.New().String()
    }
    return nil
}

type UserPreference struct {
    UserID              string    `json:"user_id" gorm:"primaryKey"`
    DefaultCategoryID   *string   `json:"default_category_id"`
    ItemsPerPage        int       `json:"items_per_page" gorm:"default:20"`
    DefaultSort         string    `json:"default_sort" gorm:"default:'saved_at'"`
    DefaultOrder        string    `json:"default_order" gorm:"default:'desc'"`
    DefaultView         string    `json:"default_view" gorm:"default:'list'"`
    Timezone           string    `json:"timezone" gorm:"default:'Asia/Tokyo'"`
    Language           string    `json:"language" gorm:"default:'ja'"`
    Theme              string    `json:"theme" gorm:"default:'light'"`
    EmailNotifications bool      `json:"email_notifications" gorm:"default:true"`
    PushNotifications  bool      `json:"push_notifications" gorm:"default:false"`
    DigestEnabled      bool      `json:"digest_enabled" gorm:"default:false"`
    DigestFrequency    *string   `json:"digest_frequency"`
    DigestTime         *string   `json:"digest_time"`
    CreatedAt          time.Time `json:"created_at"`
    UpdatedAt          time.Time `json:"updated_at"`
    
    // Relations
    User            User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
    DefaultCategory *Category `json:"default_category,omitempty" gorm:"foreignKey:DefaultCategoryID"`
}
```

### cmd/api/main.go
```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "github.com/gin-gonic/gin"
    
    "github.com/eikuma/stockle/backend/internal/config"
    "github.com/eikuma/stockle/backend/internal/database"
    "github.com/eikuma/stockle/backend/internal/controllers"
    "github.com/eikuma/stockle/backend/internal/middleware"
    "github.com/eikuma/stockle/backend/internal/repositories"
    "github.com/eikuma/stockle/backend/internal/services"
)

func main() {
    // 設定の読み込み
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    // データベース接続
    db, err := database.New(&cfg.Database)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()
    
    // マイグレーション実行
    if err := db.AutoMigrate(); err != nil {
        log.Fatalf("Failed to migrate database: %v", err)
    }
    
    // Ginの設定
    gin.SetMode(cfg.Server.Mode)
    router := gin.New()
    
    // ミドルウェアの設定
    router.Use(middleware.Logger())
    router.Use(middleware.Recovery())
    router.Use(middleware.CORS())
    router.Use(middleware.SecurityHeaders())
    
    // リポジトリとサービスの初期化
    userRepo := repositories.NewUserRepository(db.DB)
    authService := services.NewAuthService(userRepo, &cfg.JWT)
    
    // コントローラーの初期化
    authController := controllers.NewAuthController(authService)
    healthController := controllers.NewHealthController(db)
    
    // ルートの設定
    v1 := router.Group("/api/v1")
    {
        // ヘルスチェック
        v1.GET("/health", healthController.Health)
        
        // 認証
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authController.Register)
            auth.POST("/login", authController.Login)
            auth.POST("/refresh", authController.RefreshToken)
            auth.POST("/logout", authController.Logout)
        }
    }
    
    // サーバーの開始
    srv := &http.Server{
        Addr:    ":" + cfg.Server.Port,
        Handler: router,
    }
    
    // Graceful shutdown
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()
    
    log.Printf("Server started on port %s", cfg.Server.Port)
    
    // シャットダウンシグナルの待機
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    
    log.Println("Shutting down server...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
    }
    
    log.Println("Server exited")
}
```

### Makefile
```makefile
.PHONY: build run test clean dev

# 変数定義
BINARY_NAME=api
MAIN_PACKAGE=./cmd/api

# ビルド
build:
	go build -o $(BINARY_NAME) $(MAIN_PACKAGE)

# 実行
run: build
	./$(BINARY_NAME)

# 開発用サーバー（ホットリロード）
dev:
	air

# テスト実行
test:
	go test -v ./...

# テストカバレッジ
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 静的解析
lint:
	golangci-lint run

# クリーンアップ
clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html

# 依存関係の更新
deps:
	go mod tidy
	go mod download

# データベースマイグレーション
migrate-up:
	migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" up

migrate-down:
	migrate -path migrations -database "mysql://$(DB_USER):$(DB_PASSWORD)@tcp($(DB_HOST):$(DB_PORT))/$(DB_NAME)" down

# Dockerビルド
docker-build:
	docker build -t stockle-backend .

# ヘルプ
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Build and run the application"
	@echo "  dev           - Run development server with hot reload"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run static analysis"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Update dependencies"
	@echo "  migrate-up    - Run database migrations"
	@echo "  migrate-down  - Rollback database migrations"
	@echo "  docker-build  - Build Docker image"
```

## 受入条件

### 必須条件
- [ ] `make dev` でサーバーが起動する
- [ ] データベース接続が正常に動作する
- [ ] `/api/v1/health` エンドポイントが応答する
- [ ] ホットリロードが正常に動作する
- [ ] 環境変数による設定が正しく読み込まれる

### 品質条件
- [ ] ユニットテストが通る
- [ ] 静的解析でエラーが0件
- [ ] データベース接続プールが適切に設定されている
- [ ] ログが適切に出力される
- [ ] Graceful shutdownが動作する

## 推定時間
**20時間** (4-5日)

## 依存関係
- P1-001: プロジェクト基盤セットアップ

## 完了後の次ステップ
1. P1-202: 認証システム実装
2. P1-203: 記事管理API実装

## 備考
- セキュリティを最優先に実装
- パフォーマンスを意識したデータベース設計
- テストを書きながら実装を進める
- エラーハンドリングを丁寧に実装