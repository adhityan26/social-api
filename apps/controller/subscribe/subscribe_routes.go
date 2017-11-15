package subscribe

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
)

type Routes struct {
	DB           *gorm.DB
	RoutesPrefix string
}

// Handle api routing for manage user subscription
func (this *Routes) Handler(app *iris.Application) {
	controller := Controller{DB: this.DB}
	api := app.Party(this.RoutesPrefix + "/subscribe")
	{
		api.Post("/", controller.Create)
		api.Delete("/", controller.Remove)
	}
}
