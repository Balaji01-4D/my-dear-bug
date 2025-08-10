package main

import (
	"time"
	"github.com/Balaji01-4D/my-dear-bug/config"
	bug "github.com/Balaji01-4D/my-dear-bug/internals/confession"
	"github.com/Balaji01-4D/my-dear-bug/internals/tag"
	"github.com/Balaji01-4D/my-dear-bug/internals/upvote"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
	AllowOrigins:     []string{"*"}, // or "*" for dev (not for prod)
	AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE"},
	AllowHeaders:     []string{"Origin", "Content-Type"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           12 * time.Hour,
	}))
	bug.RegisterRoutes(r, db)
	upvote.RegisterRoutes(r, db)
	tag.RegisterRoutes(r, db)

	r.Run()
}
