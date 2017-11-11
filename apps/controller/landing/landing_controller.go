package landing

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
)

type Controller struct {
	DB *gorm.DB
}

func (this *Controller) Index(ctx iris.Context) {
	ctx.JSON(iris.Map{"message": "Social API v0.1!"})
}