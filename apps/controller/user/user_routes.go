package user

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
)

type Routes struct {
	DB *gorm.DB
	RoutesPrefix string
}

func (this *Routes) Handler(app *iris.Application) {
	controller := Controller{DB: this.DB}
	api := app.Party(this.RoutesPrefix + "/user")
	{
		api.Get("/", controller.Index)
		api.Post("/", controller.Create)
		api.Get("/{id}", controller.Show)
		api.Put("/{id}", controller.Update)
		api.Delete("/{id}", controller.Remove)
	}
}
