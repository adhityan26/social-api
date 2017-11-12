package connection

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
	"net/http"
	"fmt"
	"social-api/apps/models"
	"time"
)

type Controller struct {
	DB *gorm.DB
}

type connectionOutput struct {
	Friend []string `json: friend`
}

func (this *Controller) Index(ctx iris.Context) {
	ctx.StatusCode(http.StatusGone)
	ctx.JSON(iris.Map{
		"message": "Not yet implemented",
		"success": false,
	})
}

func (this *Controller) Create(ctx iris.Context) {
	param := connectionOutput{}
	ctx.ReadJSON(&param)

	var success, returnStatus = true, 200
	var message = []string{}

	if len(param.Friend) < 2 {
		returnStatus = http.StatusPreconditionRequired
		message = append(message, "Must provide 2 email")
		success = false
	}

	if (len(param.Friend) == 2) && (param.Friend[0] == param.Friend[1]) {
		returnStatus = http.StatusPreconditionFailed
		message = append(message, "Email cannot be the same")
		success = false
	}

	//if len(user.Name) == 0 {
	//	returnStatus = http.StatusPreconditionRequired
	//	message = append(message, "Name cannot be empty")
	//	success = false
	//}
	//
	if success {
		var user1 = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Friend[0]).First(&user1)
		var user2 = models.User{}
		userModel2 := this.DB.Where("email = ?", param.Friend[1]).First(&user2)

		if userModel1.RecordNotFound() {
			returnStatus = http.StatusNotFound
			message = append(message, fmt.Sprintf("Email %s not found", param.Friend[0]))
		}

		if userModel2.RecordNotFound() {
			returnStatus = http.StatusNotFound
			message = append(message, fmt.Sprintf("Email %s not found", param.Friend[1]))
			success = false
		}

		if success {
			var checkFriend1 models.Connection
			checkFriendModel1 := this.DB.
				Where("user_id = ?", user1.Id).
				Where("friend_id = ?", user2.Id).
				First(&checkFriend1)

			var checkFriend2 models.Connection
			checkFriendModel2 := this.DB.
				Where("user_id = ?", user2.Id).
				Where("friend_id = ?", user1.Id).
				First(&checkFriend2)

			if checkFriendModel1.RecordNotFound() && checkFriendModel2.RecordNotFound() {
				if success {
					tx := this.DB.Begin()
					userConnection1 := models.Connection{}
					userConnection1.UserId = user1.Id
					userConnection1.FriendId = user2.Id
					userConnection1.CreatedAt = time.Now()
					userConnection1.UpdatedAt = time.Now()
					if err := tx.Create(&userConnection1).Error; err != nil {
						returnStatus = http.StatusInternalServerError
						message = append(message, err.Error())
						success = false
					}

					if success {
						userConnection2 := models.Connection{}
						userConnection2.UserId = user2.Id
						userConnection2.FriendId = user1.Id
						userConnection2.CreatedAt = time.Now()
						userConnection2.UpdatedAt = time.Now()
						if err := tx.Create(&userConnection2).Error; err != nil {
							tx.Rollback()
							returnStatus = http.StatusInternalServerError
							message = append(message, err.Error())
							success = false
						}
					}

					if success {
						tx.Commit()
						ctx.JSON(iris.Map{
							"success": true,
						})
						return
					}
				}
			} else {
				returnStatus = http.StatusPreconditionFailed
				message = append(message, fmt.Sprintf("Email %s and %s is already friend", param.Friend[1], param.Friend[1]))
				success = false
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": message,
		"success": success,
	})
}

func (this *Controller) Remove(ctx iris.Context) {
	ctx.StatusCode(http.StatusGone)
	ctx.JSON(iris.Map{
		"message": "Not yet implemented",
		"success": false,
	})
}