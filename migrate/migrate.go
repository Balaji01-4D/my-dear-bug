package main

import (
	"github.com/Balaji01-4D/my-dear-bug/config"
	"github.com/Balaji01-4D/my-dear-bug/internals/bug"
	"github.com/Balaji01-4D/my-dear-bug/internals/upvote"
)

func init() {
	config.LoadEnvVariables()
}

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)

	db.AutoMigrate(&bug.Bug{})
	db.AutoMigrate(&upvote.Upvote{})
}
