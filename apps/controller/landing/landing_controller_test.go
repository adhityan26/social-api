package landing_test

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"social-api/apps/controller/landing"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"social-api/apps/test"
)

var server = test.Setup()

func TestLoadError(t *testing.T) {
	controller := &landing.Routes{DB: server.DB, RoutesPrefix: "/api/v1"}
	controller.Handler(server.App)

	e := httptest.New(t, server.App)

	e.GET("/").Expect().Body().NotEqual("Welcome")
	e.GET(server.RoutePrefix).Expect().Body().Contains("message")
}