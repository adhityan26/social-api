package apps

import (
	"github.com/kataras/iris"
	"social-api/apps/controller/landing"
	"social-api/apps/controller/user"
	"github.com/jinzhu/gorm"
)

type Routes struct {
	DB *gorm.DB
}

func (this *Routes) CreateApp() *iris.Application {
	app := iris.New()

	// Web
	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Welcome</b>")
	})

	// Api
	apiPrefix := "/api/v1"
	(&landing.Routes{DB: this.DB, RoutesPrefix: apiPrefix}).Handler(app)
	(&user.Routes{DB: this.DB, RoutesPrefix: apiPrefix}).Handler(app)

	return app
}
