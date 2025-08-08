package main

import (
	"github.com/Balaji01-4D/my-dear-bug/config"
	bug "github.com/Balaji01-4D/my-dear-bug/internals/confession"
	"github.com/Balaji01-4D/my-dear-bug/internals/tag"
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
	tag.RegisterRoutes(r, db)

	r.Run()
}
