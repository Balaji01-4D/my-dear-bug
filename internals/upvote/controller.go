package upvote

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"

	middleware "github.com/Balaji01-4D/shit-happens/internals/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepo(db)
	svc := NewService(repo)

	// upvote the confession
	r.POST("/confessions/:id/upvote", middleware.UpvoteRateLimitMiddleware(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		ip := c.ClientIP()

		hash := sha256.Sum256([]byte(ip))
		ipHash := hex.EncodeToString(hash[:])

		// Read or set a long-lived anonymous client id cookie
		const cookieName = "mdb_client_id"
		clientID, err := c.Cookie(cookieName)
		if err != nil || clientID == "" {
			// generate 16 random bytes hex
			buf := make([]byte, 16)
			if _, rerr := rand.Read(buf); rerr == nil {
				clientID = hex.EncodeToString(buf)
				maxAge := 60 * 60 * 24 * 365 // 1 year expiry
				// Detect HTTPS (direct or via proxy)
				secure := c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https"
				// Set cookie with SameSite=None for cross-site usage when secure
				ck := &http.Cookie{
					Name:     cookieName,
					Value:    clientID,
					Path:     "/",
					MaxAge:   maxAge,
					HttpOnly: true,
					Secure:   secure,
				}
				if secure {
					ck.SameSite = http.SameSiteNoneMode
				}
				http.SetCookie(c.Writer, ck)
			}
		}
		clientHash := ""
		if clientID != "" {
			ch := sha256.Sum256([]byte(clientID))
			clientHash = hex.EncodeToString(ch[:])
		}

		if repo.HasUpvoted(uint(id), ipHash, clientHash) {
			c.JSON(http.StatusOK, gin.H{"message": "already upvoted"})
			return
		}

		if err := svc.Upvote(uint(id), ipHash, clientHash); err != nil {
			// Possible race: re-check; if now present, treat as idempotent success
			if repo.HasUpvoted(uint(id), ipHash, clientHash) {
				c.JSON(http.StatusOK, gin.H{"message": "already upvoted"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upvote"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "upvote recorded"})
	})
}
