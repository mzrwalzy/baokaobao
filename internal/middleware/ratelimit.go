package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	burst    int
}

type visitor struct {
	tokens    float64
	lastVisit time.Time
}

var limiter *rateLimiter

func InitRateLimiter(rate, burst int) {
	limiter = &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}
	go limiter.cleanup()
}

func (rl *rateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastVisit) > time.Hour {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			tokens:    float64(rl.burst),
			lastVisit: time.Now(),
		}
		rl.visitors[ip] = v
	}
	return v
}

func (rl *rateLimiter) Allow(ip string) bool {
	v := rl.getVisitor(ip)

	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(v.lastVisit).Seconds()
	v.lastVisit = now

	v.tokens += elapsed * float64(rl.rate)
	if v.tokens > float64(rl.burst) {
		v.tokens = float64(rl.burst)
	}

	if v.tokens < 1 {
		return false
	}

	v.tokens--
	return true
}

func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if limiter == nil {
			InitRateLimiter(100, 200)
		}

		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "rate limit exceeded",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}


