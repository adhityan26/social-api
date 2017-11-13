package main

import (
	"github.com/joho/godotenv"
	"github.com/kataras/iris"
	"os"
	"social-api/apps"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"social-api/apps/models"
)

func main() {
	godotenv.Load()

	conn := os.Getenv("USERNAME_DB") + ":" + os.Getenv("PASSWORD_DB") + "@tcp(" + os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT") + ")/" + os.Getenv("DATABASE_NAME") + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(os.Getenv("DATABASE_TYPE"), conn)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	db.AutoMigrate(&models.User{}, &models.Connection{}, &models.Subscribe{})

	(&apps.Routes{DB: db}).CreateApp().Run(iris.Addr(os.Getenv("HOST") + ":" + os.Getenv("PORT")), iris.WithoutServerError(iris.ErrServerClosed))
}