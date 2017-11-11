package landing_test

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"github.com/kataras/iris"
	"social-api/apps/controller/landing"
	"github.com/joho/godotenv"
	"os"
	"github.com/jinzhu/gorm"

	_ "github.com/go-sql-driver/mysql"
)

func TestLoadError(t *testing.T) {
	godotenv.Load()

	conn := os.Getenv("USERNAME_DB") + ":" + os.Getenv("PASSWORD_DB") + "@tcp(" + os.Getenv("DATABASE_HOST") + ":" + os.Getenv("DATABASE_PORT") + ")/" + os.Getenv("DATABASE_NAME") + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(os.Getenv("DATABASE_TYPE"), conn)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := iris.New()

	controller := &landing.Routes{DB: db, RoutesPrefix: "/api/v1"}
	controller.Handler(app)

	e := httptest.New(t, app)

	e.GET("/").Expect().Body().NotEqual("Welcome")
	e.GET("/api/v1").Expect().Body().Contains("message")
}