package test

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/joho/godotenv"
	"os"
)

type Server struct {
	DB *gorm.DB
	App *iris.Application
	RoutePrefix string
}

func Setup() Server {
	server := Server{}

	godotenv.Load()

	conn := getenv("USERNAME_DB", "api") + ":" + getenv("PASSWORD_DB", "api") + "@tcp(" + getenv("DATABASE_HOST", "localhost") + ":" + getenv("DATABASE_PORT", "3336") + ")/" + getenv("DATABASE_NAME", "social-api") + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(getenv("DATABASE_TYPE", "mysql"), conn)

	if err != nil {
		panic(err)
	}

	server.DB = db
	server.App = iris.New()
	server.RoutePrefix = "/api/v1"

	return server
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}