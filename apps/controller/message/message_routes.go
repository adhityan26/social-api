package message

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
)

type Routes struct {
	DB           *gorm.DB
	RoutesPrefix string
}

// Handle api routing for sending message
func (this *Routes) Handler(app *iris.Application) {
	controller := Controller{DB: this.DB}
	api := app.Party(this.RoutesPrefix + "/message")
	{
		api.Post("/", controller.Create)
	}
}
