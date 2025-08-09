package tag

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Balaji01-4D/my-dear-bug/internals/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {
	repo := NewRepo(db)
	service := NewService(repo)

	tagRoutes := r.Group("/tags")


	tagRoutes.GET("", func(c *gin.Context) {
		tags, err := service.GetTags()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})
			return
		}

		c.JSON(http.StatusOK, tags)
	})	



	tagRoutes.POST("", func(c *gin.Context) {

		var dto struct {
			Name string `json:"name" binding:"required,min=1,max=50"`
		}

		if err := c.ShouldBindJSON(&dto); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := service.CreateTag(dto.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "cannot save tag",
			})
			return
		}

			c.Status(http.StatusOK)
	})


	tagRoutes.GET("/suggest", func (c *gin.Context)  {

		query := strings.ToLower(c.Query("query"))

		if len(query) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query too short"})
			return
		}

		tagResult, err := service.SuggestTags(query)

		if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB query failed"})

			return
		}

		c.JSON(http.StatusOK, tagResult)
		
	})

	tagRoutes.DELETE("/:id", middleware.AdminAuthMiddleware(), func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		if err := service.DeleteTags(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to delete the tag"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "tag deleted successfully"})

	})
	
}