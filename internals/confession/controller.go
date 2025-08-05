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

	Confessions := r.Group("/confessions")

	Confessions.GET("", func(c *gin.Context) {
		offest, _ := strconv.Atoi(c.Query("offset"))
		limit, _ := strconv.Atoi(c.Query("limit"))

		confessionList, err := service.List(offest, limit)

		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, confessionList)
	})

	Confessions.GET("/:id", func(c *gin.Context) {
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

	Confessions.POST("", middleware.PostRateLimitMiddleWare(), func(c *gin.Context) {
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

	Confessions.DELETE("/:id", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
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

	Confessions.GET("/language/:language", func(c *gin.Context) {
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

	Confessions.GET("/top", func(c *gin.Context) {
		offest, _ := strconv.Atoi(c.Query("offset"))
		limit, _ := strconv.Atoi(c.Query("limit"))

		confessions, err := service.GetTopConfessions(offest, limit)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch confession by rating"})
			return
		}

		c.JSON(http.StatusOK, confessions)
	})

}
