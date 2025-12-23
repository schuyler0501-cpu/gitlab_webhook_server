package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RateLimiter 限流器
type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int           // 限制数量
	window   time.Duration // 时间窗口
	logger   *zap.Logger
}

// NewRateLimiter 创建新的限流器
func NewRateLimiter(limit int, window time.Duration, logger *zap.Logger) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		logger:   logger,
	}

	// 定期清理过期记录
	go rl.cleanup()

	return rl
}

// Limit 限流中间件
func (rl *RateLimiter) Limit() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.ClientIP()

		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)

		// 清理过期请求
		requests := rl.requests[key]
		validRequests := make([]time.Time, 0)
		for _, t := range requests {
			if t.After(cutoff) {
				validRequests = append(validRequests, t)
			}
		}

		// 检查是否超过限制
		if len(validRequests) >= rl.limit {
			rl.mu.Unlock()
			rl.logger.Warn("请求被限流",
				zap.String("ip", key),
				zap.Int("count", len(validRequests)),
				zap.Int("limit", rl.limit),
			)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		// 记录本次请求
		validRequests = append(validRequests, now)
		rl.requests[key] = validRequests
		rl.mu.Unlock()

		c.Next()
	}
}

// cleanup 定期清理过期记录
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute * 5)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window * 2) // 清理两倍窗口时间的数据

		for key, requests := range rl.requests {
			validRequests := make([]time.Time, 0)
			for _, t := range requests {
				if t.After(cutoff) {
					validRequests = append(validRequests, t)
				}
			}

			if len(validRequests) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validRequests
			}
		}
		rl.mu.Unlock()
	}
}

