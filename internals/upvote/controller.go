package upvote

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepo(db)
	svc := NewService(repo)

	// upvote the confession
	r.POST("/confessions/:id/upvote", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		ip := c.ClientIP()

		hash := sha256.Sum256([]byte(ip))
		ipHash := hex.EncodeToString(hash[:])

		if repo.HasUpvoted(uint(id), ipHash) {
			return
		}

		if err := svc.Upvote(uint(id), ipHash); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "failed to upvote",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "upvote recorded"})
	})
}
