package landing

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
)

type Routes struct {
	DB *gorm.DB
	RoutesPrefix string
}

// Handle api routing for landing page
func (this *Routes) Handler(app *iris.Application) {
	controller := Controller{DB: this.DB}
	api := app.Party(this.RoutesPrefix)
	{
		api.Get("/", controller.Index)
	}
}
