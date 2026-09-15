package middleware

import (
	"net/http"
	"sync"
	"time"

	"fintech-wallet/backend/internal/response"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count       int
	windowStart time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
}

func (r *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		r.mu.Lock()

		v, exists := r.visitors[ip]

		if !exists || time.Since(v.windowStart) >= r.window {
			r.visitors[ip] = &visitor{
				count:       1,
				windowStart: time.Now(),
			}

			r.mu.Unlock()
			c.Next()
			return
		}

		if v.count >= r.limit {
			r.mu.Unlock()

			response.Error(
				c,
				http.StatusTooManyRequests,
				"Too many requests",
				"RATE_LIMIT_EXCEEDED",
			)

			c.Abort()
			return
		}

		v.count++

		r.mu.Unlock()

		c.Next()
	}
}
