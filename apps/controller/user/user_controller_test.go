package user_test

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"social-api/apps/controller/user"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"social-api/apps/test"
	"social-api/apps/models"
	"fmt"
	"github.com/kataras/iris"
)

var (
	server = test.Setup()
)

func TestCreateUser(t *testing.T) {
	defer server.DB.Close()
	controller := &user.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	controller.Handler(server.App)

	e := httptest.New(t, server.App)

	newUser := models.User{
		Name: "test",
		Email: "test@mail.com",
	}

	e.GET(server.RoutePrefix + "/user").WithForm(iris.Map{"Email": newUser.Email}).Expect().Body().Contains("User not found")

	e.GET(server.RoutePrefix + "/user/1").Expect().Status(httptest.StatusNotFound)
	e.GET(server.RoutePrefix + "/user/op").Expect().Status(httptest.StatusPreconditionFailed)

	e.POST(server.RoutePrefix + "/user").
		Expect().Status(httptest.StatusPreconditionRequired)

	e.POST(server.RoutePrefix + "/user").
		WithForm(newUser).
			Expect().Body().Contains(newUser.Email)

	var user = models.User{
		Email: newUser.Email,
	}
	server.DB.First(&user)

	if server.DB.RecordNotFound() {
		t.Error("Create user is failed")
	}

	e.GET(server.RoutePrefix + "/user").WithQueryObject(iris.Map{"Email": "test", "Name": "Te", "Status": "1"}).Expect().Body().Contains("data")

	e.GET(server.RoutePrefix + "/user/" + fmt.Sprint(user.Id)).Expect().Status(httptest.StatusOK)

	e.POST(server.RoutePrefix + "/user").
		WithForm(newUser).
		Expect().Status(httptest.StatusConflict)

	updateUser := models.User{
		Status: false,
		Name: user.Name + "-1",
	}

	e.PUT(server.RoutePrefix + "/user/0").
		WithForm(updateUser).
		Expect().Body().Contains("User not found")

	e.PUT(server.RoutePrefix + "/user/" + fmt.Sprint(user.Id)).
		WithForm(updateUser).
		Expect().Body().Contains("updated successfully")

	e.PUT(server.RoutePrefix + "/user/op").
		WithForm(updateUser).
		Expect().Status(httptest.StatusPreconditionFailed)

	user = models.User{
		Id: user.Id,
	}

	server.DB.First(&user)

	if !user.Status {
		t.Error("User is not updated")
	}

	e.DELETE(server.RoutePrefix + "/user/0").
		Expect().Body().Contains("User not found")

	e.DELETE(server.RoutePrefix + "/user/" + fmt.Sprint(user.Id)).
		Expect().Body().Contains("deleted successfully")

	e.DELETE(server.RoutePrefix + "/user/op").
		Expect().Status(httptest.StatusPreconditionFailed)
}