package main

import (
	"github.com/Balaji01-4D/my-dear-bug/config"
	"github.com/Balaji01-4D/my-dear-bug/internals/bug"
	"github.com/Balaji01-4D/my-dear-bug/internals/upvote"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	bug.RegisterRoutes(r, db)
	upvote.RegisterRoutes(r, db)

	r.Run()
}
