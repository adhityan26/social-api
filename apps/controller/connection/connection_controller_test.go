package connection_test

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"social-api/apps/controller/user"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"social-api/apps/test"
	"social-api/apps/models"
	"fmt"
	"github.com/kataras/iris"
	"social-api/apps/controller/connection"
)

var (
	server = test.Setup()
)

func TestCreateConnection(t *testing.T) {
	defer server.DB.Close()
	controller := &connection.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	controller.Handler(server.App)

	userController := &user.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	userController.Handler(server.App)

	e := httptest.New(t, server.App)

	//init test user
	newUser := models.User{
		Name: "test",
		Email: "test@mail.com",
	}

	for i := 0; i < 10; i++ {
		dataUser := iris.Map{
			"name": newUser.Name + "_" + fmt.Sprint(i),
			"email": "mail_" + fmt.Sprint(i) + "_" + newUser.Email,
		}
		e.POST(server.RoutePrefix + "/user").
			WithJSON(dataUser).
			Expect().Status(httptest.StatusOK)
	}

	var userList []models.User
	for i := 0; i < 10; i++ {
		user := models.User{}

		server.DB.Where("email = ?", "mail_" + fmt.Sprint(i) + "_" + newUser.Email).First(&user)

		userList = append(userList, user)
	}

	e.POST(server.RoutePrefix + "/connection").
		WithJSON(iris.Map{"friends": [2]string{userList[0].Email, userList[0].Email}}).Expect().Status(httptest.StatusPreconditionFailed)

	e.POST(server.RoutePrefix + "/connection").
		Expect().Status(httptest.StatusPreconditionRequired)

	e.POST(server.RoutePrefix + "/connection").
		WithJSON(iris.Map{"friends": [1]string{userList[0].Email}}).Expect().Status(httptest.StatusPreconditionRequired)

	for i := 0; i < 5; i++ {
		e.POST(server.RoutePrefix + "/connection").
			WithJSON(iris.Map{"friends": [2]string{userList[0].Email, userList[i + 1].Email}}).
				Expect().Body().Contains("\"success\":true")
	}
	var userConnection []models.Connection

	server.DB.Where("user_id = ?", userList[0].Id).Find(&userConnection)
	if len(userConnection) != 5 {
		t.Error("Not all user is connected")
	}

	for i := 9; i > 5; i-- {
		e.POST(server.RoutePrefix + "/connection").
			WithJSON(iris.Map{"friends": [2]string{userList[1].Email, userList[i].Email}}).
			Expect().Body().Contains("\"success\":true")
	}

	server.DB.Where("user_id = ?", userList[1].Id).Find(&userConnection)
	if len(userConnection) != 5 {
		t.Error("Not all user is connected")
	}

	e.POST(server.RoutePrefix + "/connection").
		WithJSON(iris.Map{"friends": [2]string{userList[1].Email, userList[7].Email}}).
			Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection").
		WithJSON(iris.Map{"friends": [2]string{"a@a.com", "a@b.com"}}).
			Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection/show").
		WithJSON(iris.Map{"email": userList[0].Email}).
			Expect().Body().Contains("\"count\":5")

	e.POST(server.RoutePrefix + "/connection/show").
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection/show").
		WithJSON(iris.Map{"email": userList[1].Email}).
			Expect().Body().Contains(userList[0].Email)

	e.POST(server.RoutePrefix + "/connection/common").
		WithJSON(iris.Map{"friends": [1]string{userList[0].Email}}).
			Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection/common").
		WithJSON(iris.Map{"friends": [2]string{"b@b.com", "a@a.com"}}).
			Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection/common").
		WithJSON(iris.Map{"friends": [2]string{"a@a.com", "a@a.com"}}).
			Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection/common").
		WithJSON(iris.Map{"friends": [2]string{userList[0].Email, userList[9].Email}}).
			Expect().Body().Contains(userList[1].Email)

	e.POST(server.RoutePrefix + "/connection/common").
		WithJSON(iris.Map{"friends": [2]string{userList[0].Email, userList[2].Email}}).
			Expect().Body().Contains("\"success\":false")

	//remove test data
	for _, us := range userList {
		server.DB.Delete(models.User{}, us.Id)
		server.DB.Where("user_id = ?", us.Id).Delete(models.Connection{})
	}
}