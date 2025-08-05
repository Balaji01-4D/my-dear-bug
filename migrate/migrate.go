package main

import (
	"github.com/Balaji01-4D/my-dear-bug/config"
	confession "github.com/Balaji01-4D/my-dear-bug/internals/confession"
	"github.com/Balaji01-4D/my-dear-bug/internals/upvote"
)

// func init() {
// 	config.LoadEnvVariables()
// }

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)

	db.AutoMigrate(&confession.Confession{})
	db.AutoMigrate(&upvote.Upvote{})
}
