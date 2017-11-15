package block_test

import (
	"testing"
	"github.com/kataras/iris/httptest"
	"social-api/apps/controller/user"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"social-api/apps/test"
	"social-api/apps/models"
	"fmt"
	"github.com/kataras/iris"
	"social-api/apps/controller/block"
	"social-api/apps/controller/connection"
	"social-api/apps/controller/subscribe"
)

var (
	server = test.Create()
)

func TestBlockUser(t *testing.T) {
	defer server.DB.Close()
	controller := &block.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	controller.Handler(server.App)

	connectionController := &connection.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	connectionController.Handler(server.App)

	subscribeController := &subscribe.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	subscribeController.Handler(server.App)

	userController := &user.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	userController.Handler(server.App)

	e := httptest.New(t, server.App)

	//init test user
	newUser := models.User{
		Name:  "test",
		Email: "test@mail.com",
	}

	for i := 0; i < 10; i++ {
		dataUser := iris.Map{
			"name":  newUser.Name + "_" + fmt.Sprint(i),
			"email": "block_mail_" + fmt.Sprint(i) + "_" + newUser.Email,
		}
		e.POST(server.RoutePrefix + "/user").
			WithJSON(dataUser).
			Expect().Status(httptest.StatusOK)
	}

	var userList []models.User
	for i := 0; i < 10; i++ {
		user := models.User{}

		server.DB.Where("email = ?", "block_mail_"+fmt.Sprint(i)+"_"+newUser.Email).First(&user)

		userList = append(userList, user)
	}

	e.POST(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": userList[0].Email, "target": userList[0].Email}).
		Expect().Status(httptest.StatusPreconditionFailed)

	e.POST(server.RoutePrefix + "/block").
		Expect().Status(httptest.StatusPreconditionRequired)

	for i := 0; i < 5; i++ {
		e.POST(server.RoutePrefix + "/block").
			WithJSON(iris.Map{"requestor": userList[0].Email, "target": userList[i+1].Email}).
			Expect().Body().Contains("\"success\":true")
	}
	var userBlock []models.Block

	server.DB.Where("requestor_id = ?", userList[0].Id).Find(&userBlock)
	if len(userBlock) != 5 {
		t.Error("Not all user is blocked")
	}

	for i := 9; i > 5; i-- {
		e.POST(server.RoutePrefix + "/block").
			WithJSON(iris.Map{"requestor": userList[1].Email, "target": userList[i].Email}).
			Expect().Body().Contains("\"success\":true")
	}

	server.DB.Where("requestor_id = ?", userList[1].Id).Find(&userBlock)
	if len(userBlock) != 4 {
		t.Error("Not all user is blocked")
	}

	e.POST(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": userList[0].Email, "target": userList[1].Email}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": "aaa", "target": "aaa"}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": "a@a.com", "target": "a@b.com"}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection").
		WithJSON(iris.Map{"friends": [2]string{userList[0].Email, userList[1].Email}}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection").
		WithJSON(iris.Map{"friends": [2]string{userList[2].Email, userList[0].Email}}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/connection").
		WithJSON(iris.Map{"friends": [2]string{userList[0].Email, userList[8].Email}}).
		Expect().Body().Contains("\"success\":true")

	e.POST(server.RoutePrefix + "/subscribe").
		WithJSON(iris.Map{"requestor": userList[1].Email, "target": userList[0].Email}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/subscribe").
		WithJSON(iris.Map{"requestor": userList[2].Email, "target": userList[0].Email}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/subscribe").
		WithJSON(iris.Map{"requestor": userList[0].Email, "target": userList[1].Email}).
		Expect().Body().Contains("\"success\":true")

	e.DELETE(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": userList[0].Email, "target": userList[1].Email}).
		Expect().Body().Contains("\"success\":true")

	e.DELETE(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": "a@a.com", "target": "a@a.com"}).
		Expect().Body().Contains("\"success\":false")

	e.DELETE(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": "a@a.com", "target": "b@b.com"}).
		Expect().Body().Contains("\"success\":false")

	e.DELETE(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": "aaa", "target": "aaa"}).
		Expect().Body().Contains("\"success\":false")

	e.DELETE(server.RoutePrefix + "/block").
		Expect().Body().Contains("\"success\":false")

	//remove test data
	for _, us := range userList {
		server.DB.Delete(models.User{}, us.Id)
		server.DB.Where("requestor_id = ?", us.Id).Delete(models.Block{})
		server.DB.Where("requestor_id = ?", us.Id).Delete(models.Subscribe{})
		server.DB.Where("user_id = ?", us.Id).Delete(models.Connection{})
	}
}
