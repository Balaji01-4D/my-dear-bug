package main

import (
	"time"

	"github.com/Balaji01-4D/shit-happens/config"
	bug "github.com/Balaji01-4D/shit-happens/internals/confession"
	"github.com/Balaji01-4D/shit-happens/internals/tag"
	"github.com/Balaji01-4D/shit-happens/internals/upvote"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://shit-happens.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	bug.RegisterRoutes(r, db)
	upvote.RegisterRoutes(r, db)
	tag.RegisterRoutes(r, db)

	r.Run()
}
