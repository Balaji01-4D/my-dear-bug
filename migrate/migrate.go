package main

import (
	"github.com/Balaji01-4D/shit-happens/config"
	confession "github.com/Balaji01-4D/shit-happens/internals/confession"
	"github.com/Balaji01-4D/shit-happens/internals/tag"
	"github.com/Balaji01-4D/shit-happens/internals/upvote"
)

func init() {
	config.LoadEnvVariables()
}

func main() {
	cfg := config.Load()
	db := config.InitDB(cfg)

	db.AutoMigrate(&confession.Confession{})
	db.AutoMigrate(&upvote.Upvote{})
	db.AutoMigrate(&tag.Tag{})
}
