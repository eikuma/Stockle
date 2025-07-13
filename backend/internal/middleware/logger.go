package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/private/Stockle/backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func InitLogger(cfg *config.Config) error {
	var zapConfig zap.Config
	
	if cfg.Server.Environment == "production" {
		zapConfig = zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(getLogLevel(cfg.Log.Level))
	} else {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(getLogLevel(cfg.Log.Level))
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	if cfg.Log.Format == "json" {
		zapConfig.Encoding = "json"
	} else {
		zapConfig.Encoding = "console"
	}
	
	var err error
	logger, err = zapConfig.Build()
	if err != nil {
		return err
	}
	
	return nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		if logger != nil {
			logger.Info("HTTP Request",
				zap.String("client_ip", param.ClientIP),
				zap.String("method", param.Method),
				zap.String("path", param.Path),
				zap.String("protocol", param.Request.Proto),
				zap.Int("status_code", param.StatusCode),
				zap.Duration("latency", param.Latency),
				zap.String("user_agent", param.Request.UserAgent()),
				zap.Int("body_size", param.BodySize),
			)
		}
		return ""
	})
}

func GetLogger() *zap.Logger {
	return logger
}