package main

import (
	"github.com/kataras/iris"
	"os"
	"social-api/apps"

	"social-api/apps/models"
	"social-api/apps/helper"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := helper.SetupDB()

	// close db connection after application is terminated
	defer db.Close()

	// migrate database using models
	db.AutoMigrate(&models.User{}, &models.Connection{}, &models.Subscribe{}, &models.Block{}, &models.Message{})

	// serve application
	(&apps.Routes{DB: db}).CreateApp().Run(iris.Addr(os.Getenv("HOST")+":"+os.Getenv("PORT")), iris.WithoutServerError(iris.ErrServerClosed))
}
