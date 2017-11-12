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

	newUser := iris.Map{
		"name": "test",
		"email": "test@mail.com",
	}

	e.GET(server.RoutePrefix + "/user").
		WithQueryObject(iris.Map{"email": newUser["email"]}).
			Expect().Body().Contains("User not found")

	e.GET(server.RoutePrefix + "/user/1").Expect().Status(httptest.StatusNotFound)
	e.GET(server.RoutePrefix + "/user/op").Expect().Status(httptest.StatusPreconditionFailed)

	e.POST(server.RoutePrefix + "/user").
		Expect().Status(httptest.StatusPreconditionRequired)

	e.POST(server.RoutePrefix + "/user").
		WithJSON(newUser).
			Expect().Body().Contains(newUser["email"].(string))

	var user = models.User{
		Email: newUser["email"].(string),
	}
	server.DB.Where("email = ?", user.Email).First(&user)

	if server.DB.RecordNotFound() {
		t.Error("Create user is failed")
	}

	e.GET(server.RoutePrefix + "/user").
		WithQueryObject(iris.Map{"email": "test", "name": "Te", "status": "1"}).
			Expect().Body().Contains("data")

	e.GET(server.RoutePrefix + "/user/" + fmt.Sprint(user.Id)).
		Expect().Status(httptest.StatusOK)

	e.POST(server.RoutePrefix + "/user").
		WithJSON(newUser).
		Expect().Status(httptest.StatusConflict)

	updateUser := iris.Map{
		"status": "0",
		"name": user.Name + "-1",
	}

	e.PUT(server.RoutePrefix + "/user/0").
		WithJSON(updateUser).
		Expect().Body().Contains("User not found")

	e.PUT(server.RoutePrefix + "/user/" + fmt.Sprint(user.Id)).
		WithJSON(updateUser).
		Expect().Body().Contains("updated successfully")

	e.PUT(server.RoutePrefix + "/user/op").
		WithJSON(updateUser).
		Expect().Status(httptest.StatusPreconditionFailed)

	user = models.User{
		Id: user.Id,
	}

	server.DB.Where("email = ?", newUser["email"]).First(&user)

	if user.Status {
		t.Error("User is not updated")
	}

	e.DELETE(server.RoutePrefix + "/user/0").
		Expect().Body().Contains("User not found")

	e.DELETE(server.RoutePrefix + "/user/" + fmt.Sprint(user.Id)).
		Expect().Body().Contains("deleted successfully")

	e.DELETE(server.RoutePrefix + "/user/op").
		Expect().Status(httptest.StatusPreconditionFailed)
}