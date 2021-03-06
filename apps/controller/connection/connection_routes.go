package connection

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
)

type Routes struct {
	DB           *gorm.DB
	RoutesPrefix string
}

// Handle api routing for manage friends connection
func (this *Routes) Handler(app *iris.Application) {
	controller := Controller{DB: this.DB}
	api := app.Party(this.RoutesPrefix + "/connection")
	{
		api.Post("/show", controller.Index)
		api.Post("/", controller.Create)
		api.Delete("/", controller.Remove)
		api.Post("/common", controller.Common)
	}
}
