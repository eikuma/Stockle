package database

import (
	"fmt"
	"log"
	"time"

	"github.com/private/Stockle/backend/internal/config"
	"github.com/private/Stockle/backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) error {
	var err error
	
	dsn := cfg.Database.DSN()
	
	// GORM logger configuration
	var gormLogger logger.Interface
	if cfg.Server.Environment == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}
	
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Get underlying sql.DB instance
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	// Set connection pool settings
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.Database.ConnMaxLifetime)
	
	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}
	
	log.Println("Database connected successfully")
	return nil
}

func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	err := DB.AutoMigrate(
		&models.User{},
		&models.UserSession{},
		&models.UserPreference{},
	)
	
	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}
	
	log.Println("Database migration completed successfully")
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close database connection: %w", err)
	}
	
	log.Println("Database connection closed")
	return nil
}

func GetDB() *gorm.DB {
	return DB
}

func IsHealthy() error {
	if DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}
	
	return nil
}