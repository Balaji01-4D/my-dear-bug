package confession

import (
	"net/http"
	"strconv"

	"github.com/Balaji01-4D/my-dear-bug/internals/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepo(db)
	service := NewService(repo)

	confessionRoutes := r.Group("/confessions")

	confessionRoutes.GET("", func(c *gin.Context) {
		offest, _ := strconv.Atoi(c.Query("offset"))
		limit, _ := strconv.Atoi(c.Query("limit"))

		confessionList, err := service.List(offest, limit)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, confessionList)
	})

	confessionRoutes.GET("/:id", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		confession, err := service.Get(uint(id))

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "not found",
			})
			return
		}

		c.JSON(http.StatusOK, confession)
	})

	confessionRoutes.POST("", middleware.PostRateLimitMiddleWare(), func(c *gin.Context) {
		var dto ConfessionRequest

		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		confession, err := service.Create(dto)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, confession)
	})

	confessionRoutes.DELETE("/:id", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		if err := service.Delete(uint(id)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "successfully deleted",
		})
	})

	confessionRoutes.GET("/language/:language", func(c *gin.Context) {
		language := c.Param("language")
		offest, _ := strconv.Atoi(c.Query("offset"))
		limit, _ := strconv.Atoi(c.Query("limit"))

		confessions, err := service.GetByLanguage(language, offest, limit)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to fetch",
			})
			return
		}

		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/top", func(c *gin.Context) {
		offest, _ := strconv.Atoi(c.Query("offset"))
		limit, _ := strconv.Atoi(c.Query("limit"))

		confessions, err := service.GetTopConfessions(offest, limit)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch confession by rating"})
			return
		}

		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/trending/weekly", func(c *gin.Context) {
		offset, _ := strconv.Atoi(c.Query("offset"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		confessions, err := service.TrendingWeekly(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/trending/monthly", func(c *gin.Context) {
		offset, _ := strconv.Atoi(c.Query("offset"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		confessions, err := service.TrendingMonthly(offset, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, confessions)
	})

	confessionRoutes.GET("/hall-of-fame", func(c *gin.Context) {
		offset, _ := strconv.Atoi(c.Query("offset"))
		limit, _ := strconv.Atoi(c.Query("limit"))
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
			return
		}
		c.JSON(http.StatusOK, cfs)
	})

}
