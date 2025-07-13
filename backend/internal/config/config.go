package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Log      LogConfig      `mapstructure:"log"`
	AI       AIConfig       `mapstructure:"ai"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"port"`
	Environment  string        `mapstructure:"environment"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
	CORS         CORSConfig    `mapstructure:"cors"`
}

type CORSConfig struct {
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	AllowedMethods []string `mapstructure:"allowed_methods"`
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	User            string        `mapstructure:"user"`
	Password        string        `mapstructure:"password"`
	Database        string        `mapstructure:"database"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Database)
}

type JWTConfig struct {
	AccessSecret   string        `mapstructure:"access_secret"`
	RefreshSecret  string        `mapstructure:"refresh_secret"`
	AccessExpiry   time.Duration `mapstructure:"access_expiry"`
	RefreshExpiry  time.Duration `mapstructure:"refresh_expiry"`
	Issuer         string        `mapstructure:"issuer"`
	CookieDomain   string        `mapstructure:"cookie_domain"`
	CookieSecure   bool          `mapstructure:"cookie_secure"`
	CookieHTTPOnly bool          `mapstructure:"cookie_http_only"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"` // json or console
}

// AIConfig is defined in ai_config.go to avoid duplication

var cfg *Config

func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	
	viper.AutomaticEnv()
	
	// デフォルト値を設定
	setDefaults()
	
	// 環境変数のマッピング
	bindEnvVars()
	
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}
	
	cfg = &config
	return cfg, nil
}

func Get() *Config {
	return cfg
}

func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.environment", "development")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")
	viper.SetDefault("server.cors.allowed_origins", []string{"http://localhost:3000"})
	viper.SetDefault("server.cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"})
	viper.SetDefault("server.cors.allowed_headers", []string{"Origin", "Content-Type", "Authorization"})
	
	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 3306)
	viper.SetDefault("database.user", "root")
	viper.SetDefault("database.password", "")
	viper.SetDefault("database.database", "stockle_db")
	viper.SetDefault("database.max_open_conns", 25)
	viper.SetDefault("database.max_idle_conns", 25)
	viper.SetDefault("database.conn_max_lifetime", "5m")
	
	// JWT defaults
	viper.SetDefault("jwt.access_expiry", "15m")
	viper.SetDefault("jwt.refresh_expiry", "7d")
	viper.SetDefault("jwt.issuer", "stockle-api")
	viper.SetDefault("jwt.cookie_domain", "localhost")
	viper.SetDefault("jwt.cookie_secure", false)
	viper.SetDefault("jwt.cookie_http_only", true)
	
	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.format", "console")
	
	// AI defaults
	viper.SetDefault("ai.request_timeout", "30s")
	viper.SetDefault("ai.max_retries", 3)
	viper.SetDefault("ai.retry_delay", "1s")
	viper.SetDefault("ai.rate_limit_per_min", 60)
}

func bindEnvVars() {
	// Server
	viper.BindEnv("server.port", "PORT")
	viper.BindEnv("server.environment", "ENV")
	
	// Database
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.database", "DB_NAME")
	
	// JWT
	viper.BindEnv("jwt.access_secret", "JWT_ACCESS_SECRET")
	viper.BindEnv("jwt.refresh_secret", "JWT_REFRESH_SECRET")
	
	// AI
	viper.BindEnv("ai.groq_api_key", "GROQ_API_KEY")
	viper.BindEnv("ai.anthropic_api_key", "ANTHROPIC_API_KEY")
}

func validateConfig(config *Config) error {
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}
	
	if config.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	
	if config.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT access secret is required")
	}
	
	if config.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT refresh secret is required")
	}
	
	return nil
}

func IsDevelopment() bool {
	if cfg == nil {
		return true
	}
	return cfg.Server.Environment == "development"
}

func IsProduction() bool {
	if cfg == nil {
		return false
	}
	return cfg.Server.Environment == "production"
}