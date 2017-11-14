package test

import (
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"social-api/apps/helper"
	"github.com/joho/godotenv"
)

type Server struct {
	DB          *gorm.DB
	App         *iris.Application
	RoutePrefix string
}

func Create() Server {
	godotenv.Load()

	server := Server{}

	server.DB = helper.SetupDB()
	server.App = iris.New()
	server.RoutePrefix = "/api/v1"

	return server
}
