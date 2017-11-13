package message_test

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
	"social-api/apps/controller/message"
)

var (
	server = test.Setup()
)

func TestBlockUser(t *testing.T) {
	defer server.DB.Close()
	controller := &message.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	controller.Handler(server.App)

	blockController := &block.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	blockController.Handler(server.App)

	connectionController := &connection.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	connectionController.Handler(server.App)

	subscribeController := &subscribe.Routes{DB: server.DB, RoutesPrefix: server.RoutePrefix}
	subscribeController.Handler(server.App)

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

	e.POST(server.RoutePrefix + "/block").
		WithJSON(iris.Map{"requestor": userList[0].Email, "target": userList[1].Email}).
			Expect().Body().Contains("\"success\":true")

	for i := 1; i < 5; i++ {
		e.POST(server.RoutePrefix + "/connection").
			WithJSON(iris.Map{"friends": [2]string{userList[0].Email, userList[i + 1].Email}}).
			Expect().Body().Contains("\"success\":true")
	}

	for i := 9; i > 5; i-- {
		e.POST(server.RoutePrefix + "/subscribe").
			WithJSON(iris.Map{"requestor": userList[1].Email, "target": userList[i].Email}).
			Expect().Body().Contains("\"success\":true")
	}

	e.POST(server.RoutePrefix + "/message").
		WithJSON(iris.Map{"sender": userList[0].Email, "text": "Hello"}).
		Expect().Body().Contains("\"success\":true")

	e.POST(server.RoutePrefix + "/message").
		WithJSON(iris.Map{"sender": userList[0].Email, "text": "Hello adhityanugraha@gmail.com, farahbellanadia@gmail.com, mail_0_test@mail.com, mail_1_test@mail.com "}).
		Expect().Body().Contains("mail_2_test@mail.com")

	e.POST(server.RoutePrefix + "/message").
		WithJSON(iris.Map{"sender": userList[6].Email, "text": "Hello mail_0_test@mail.com"}).
		Expect().Body().Contains("mail_1_test@mail.com")

	e.POST(server.RoutePrefix + "/message").
		WithJSON(iris.Map{"sender": userList[1].Email + "a", "text": "Hello"}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/message").
		WithJSON(iris.Map{"sender": userList[1].Email, "text": "Hello"}).
		Expect().Body().Contains("\"success\":false")

	e.POST(server.RoutePrefix + "/message").
		Expect().Body().Contains("\"success\":false")

	//remove test data
	for _, us := range userList {
		server.DB.Delete(models.User{}, us.Id)
		server.DB.Where("requestor_id = ?", us.Id).Delete(models.Block{})
		server.DB.Where("requestor_id = ?", us.Id).Delete(models.Subscribe{})
		server.DB.Where("user_id = ?", us.Id).Delete(models.Connection{})
		server.DB.Where("sender_id = ?", us.Id).Delete(models.Message{})
	}
}