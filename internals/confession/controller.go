package confession

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Balaji01-4D/shit-happens/internals/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	defaultLimit = 10
	maxLimit     = 100
)

func parsePagination(c *gin.Context) (offset, limit int) {
	offset, _ = strconv.Atoi(c.Query("offset"))
	limit, _ = strconv.Atoi(c.Query("limit"))
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return
}

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepo(db)
	service := NewService(repo)

	confessionRoutes := r.Group("/confessions")

	confessionRoutes.GET("", func(c *gin.Context) {
		offset, limit := parsePagination(c)
		list, err := service.List(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
			return
		}
		c.JSON(http.StatusOK, list)
	})

	confessionRoutes.GET("/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		confession, err := service.Get(uint(id))
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
			return
		}
		c.JSON(http.StatusOK, confession)
	})

	confessionRoutes.POST("", middleware.PostRateLimitMiddleWare(), func(c *gin.Context) {
		var dto ConfessionRequest
		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		confession, err := service.Create(dto)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create"})
			return
		}
		c.JSON(http.StatusCreated, confession)
	})

	confessionRoutes.DELETE("/:id", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := service.Delete(uint(id)); err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "successfully deleted"})
	})

	confessionRoutes.GET("/language/:language", func(c *gin.Context) {
		language := strings.TrimSpace(c.Param("language"))
		if language == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "language required"})
			return
		}
		offset, limit := parsePagination(c)
		confessions, err := service.GetByLanguage(language, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
			return
		}
		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/top", func(c *gin.Context) {
		offset, limit := parsePagination(c)
		confessions, err := service.GetTopConfessions(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch"})
			return
		}
		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/trending/weekly", func(c *gin.Context) {
		offset, limit := parsePagination(c)
		confessions, err := service.TrendingWeekly(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/trending/monthly", func(c *gin.Context) {
		offset, limit := parsePagination(c)
		confessions, err := service.TrendingMonthly(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/hall-of-fame", func(c *gin.Context) {
		offset, limit := parsePagination(c)
		confessions, err := service.HallOfFame(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/random", func(c *gin.Context) {
		cfs, err := service.Random()
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, cfs)
	})

	confessionRoutes.GET("/search", func(c *gin.Context) {
		q := strings.TrimSpace(c.Query("q"))

		language := strings.TrimSpace(c.Query("language"))
		tag := strings.TrimSpace(c.Query("tag"))

		if q == "" && language == "" && tag == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "query parameters required"})
			return
		}

		offset, limit := parsePagination(c)
		results, err := service.Search(q, language, tag, offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, results)
	})
}
