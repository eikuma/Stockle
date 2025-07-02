package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if logger != nil {
			logger.Error("Panic recovered",
				zap.Any("error", recovered),
				zap.String("path", c.Request.URL.Path),
			)
		}
		
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
	})
}