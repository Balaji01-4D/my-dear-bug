package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var visitors = make(map[string]*Visitor)
var mu sync.Mutex

// clear visitors for every ten minutes
func init() {
	go func() {
		for {
			time.Sleep(time.Minute * 10)
			mu.Lock()

			for ip, v := range visitors {
				if time.Since(v.lastSeen) > time.Minute*10 {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	visitor, exists := visitors[ip]

	if !exists {
		limiter := rate.NewLimiter(10.0/3600.0, 3) // 10 posts per hour
		visitors[ip] = &Visitor{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	visitor.lastSeen = time.Now()
	return visitor.limiter
}

func PostRateLimitMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		limiter := getVisitor(ip)

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests - slow down",
			})
			return
		}

		c.Next()

	}
}
