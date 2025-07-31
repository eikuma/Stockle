package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	mu      sync.RWMutex
	clients map[string]*clientInfo
}

type clientInfo struct {
	lastSeen time.Time
	count    int
}

func NewRateLimiter() *rateLimiter {
	rl := &rateLimiter{
		clients: make(map[string]*clientInfo),
	}

	// クリーンアップゴルーチン
	go rl.cleanup()

	return rl
}

func (rl *rateLimiter) RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()

		if client, exists := rl.clients[ip]; exists {
			if now.Sub(client.lastSeen) > window {
				client.count = 1
				client.lastSeen = now
			} else {
				client.count++
				if client.count > maxRequests {
					c.JSON(http.StatusTooManyRequests, gin.H{
						"error": "Rate limit exceeded",
						"code":  "RATE_LIMIT_EXCEEDED",
					})
					c.Abort()
					return
				}
			}
		} else {
			rl.clients[ip] = &clientInfo{
				lastSeen: now,
				count:    1,
			}
		}

		c.Next()
	}
}

func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)

		rl.mu.Lock()
		now := time.Now()
		for ip, client := range rl.clients {
			if now.Sub(client.lastSeen) > time.Hour {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}