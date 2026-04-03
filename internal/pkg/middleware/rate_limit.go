package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateBucket struct {
	tokens     int
	lastRefill time.Time
}

func RateLimitMiddleware(rps int, burst int) gin.HandlerFunc {
	if rps <= 0 || burst <= 0 {
		return func(c *gin.Context) { c.Next() }
	}

	var (
		mu      sync.Mutex
		buckets = make(map[string]*rateBucket)
	)

	refillInterval := time.Second / time.Duration(rps)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		bucket, ok := buckets[ip]
		if !ok {
			bucket = &rateBucket{tokens: burst, lastRefill: now}
			buckets[ip] = bucket
		}

		elapsed := now.Sub(bucket.lastRefill)
		if elapsed > 0 {
			newTokens := int(elapsed / refillInterval)
			if newTokens > 0 {
				bucket.tokens += newTokens
				if bucket.tokens > burst {
					bucket.tokens = burst
				}
				bucket.lastRefill = bucket.lastRefill.Add(time.Duration(newTokens) * refillInterval)
			}
		}

		if bucket.tokens <= 0 {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			c.Abort()
			return
		}

		bucket.tokens--
		mu.Unlock()

		c.Next()
	}
}
