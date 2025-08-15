package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// Lightweight per-client upvote rate limiter.
// Key is derived from cookie if present, else IP.
type upvoteVisitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	upvoteVisitors = make(map[string]*upvoteVisitor)
	upvoteMu       sync.Mutex
)

func upvoteKey(c *gin.Context) string {
	if id, err := c.Cookie("mdb_client_id"); err == nil && id != "" {
		h := sha256.Sum256([]byte(id))
		return hex.EncodeToString(h[:])
	}
	return c.ClientIP()
}

func getUpvoteLimiter(key string) *rate.Limiter {
	upvoteMu.Lock()
	defer upvoteMu.Unlock()
	v, ok := upvoteVisitors[key]
	if !ok {
		// Allow 1 upvote per 10 seconds on average with burst of 3
		lim := rate.NewLimiter(rate.Every(10*time.Second), 3)
		upvoteVisitors[key] = &upvoteVisitor{limiter: lim, lastSeen: time.Now()}
		return lim
	}
	v.lastSeen = time.Now()
	return v.limiter
}

func init() {
	// cleanup goroutine
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			upvoteMu.Lock()
			for k, v := range upvoteVisitors {
				if time.Since(v.lastSeen) > 10*time.Minute {
					delete(upvoteVisitors, k)
				}
			}
			upvoteMu.Unlock()
		}
	}()
}

func UpvoteRateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := upvoteKey(c)
		if !getUpvoteLimiter(key).Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Too many upvotes, slow down"})
			return
		}
		c.Next()
	}
}
