package config

import (
	"os"
	"log"

	"github.com/joho/godotenv"
	"github.com/Balaji01-4D/my-dear-bug/internals/bug"
	"github.com/Balaji01-4D/my-dear-bug/internals/upvote"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DBUrl         string
	ServerAddress string
}


func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading the .env files", err)
		
	}
}

func init() {
	LoadEnvVariables()
}

func Load() *Config {
	
	return &Config{
		DBUrl:         os.Getenv("DATABASE_URL"),
		ServerAddress: os.Getenv("SERVER_ADDR"),
	}
}

func InitDB(cfg *Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
