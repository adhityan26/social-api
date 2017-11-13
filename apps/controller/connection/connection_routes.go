package connection

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
	api := app.Party(this.RoutesPrefix + "/connection")
	{
		api.Post("/show", controller.Index)
		api.Post("/", controller.Create)
		api.Post("/common", controller.Common)
		api.Delete("/remove", controller.Remove)
	}
}
