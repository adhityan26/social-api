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
	Friends []string `json: friends`
}

type friendList struct {
	Email string `json: email`
}

func (this *Controller) Index(ctx iris.Context) {
	param := friendList{}
	ctx.ReadJSON(&param)

	var success, returnStatus = true, 200
	var message = []string{}

	if len(param.Email) == 0 {
		returnStatus = http.StatusPreconditionRequired
		message = append(message, "Email should not be empty")
		success = false
	}

	if success {
		var user = models.User{}
		userModel := this.DB.Where("email = ?", param.Email).
			Preload("Friends").
				Preload("Friends.UserDetail").First(&user)
		if userModel.RecordNotFound() {
			returnStatus = http.StatusPreconditionFailed
			message = append(message, "User not found")
			success = false
		} else {
			listEmail := []string{}
			for _, friend := range user.Friends {
				listEmail = append(listEmail, friend.UserDetail.Email)
			}

			ctx.JSON(iris.Map{
				"friends": listEmail,
				"count": len(user.Friends),
				"success": success,
			})
			return
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": message,
		"success": success,
	})
}

func (this *Controller) Create(ctx iris.Context) {
	param := connectionOutput{}
	ctx.ReadJSON(&param)

	var success, returnStatus = true, 200
	var message = []string{}

	if len(param.Friends) < 2 {
		returnStatus = http.StatusPreconditionRequired
		message = append(message, "Must provide 2 email")
		success = false
	}

	if (len(param.Friends) == 2) && (param.Friends[0] == param.Friends[1]) {
		returnStatus = http.StatusPreconditionFailed
		message = append(message, "Email cannot be the same")
		success = false
	}

	if success {
		var user1 = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Friends[0]).First(&user1)
		var user2 = models.User{}
		userModel2 := this.DB.Where("email = ?", param.Friends[1]).First(&user2)

		if userModel1.RecordNotFound() {
			returnStatus = http.StatusNotFound
			message = append(message, fmt.Sprintf("Email %s not found", param.Friends[0]))
		}

		if userModel2.RecordNotFound() {
			returnStatus = http.StatusNotFound
			message = append(message, fmt.Sprintf("Email %s not found", param.Friends[1]))
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
				message = append(message, fmt.Sprintf("Email %s and %s is already friend", param.Friends[1], param.Friends[1]))
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

func (this *Controller) Common(ctx iris.Context) {
	param := connectionOutput{}
	ctx.ReadJSON(&param)

	var success, returnStatus = true, 200
	var message = []string{}

	if len(param.Friends) < 2 {
		returnStatus = http.StatusPreconditionRequired
		message = append(message, "Must provide 2 email")
		success = false
	}

	if (len(param.Friends) == 2) && (param.Friends[0] == param.Friends[1]) {
		returnStatus = http.StatusPreconditionFailed
		message = append(message, "Email cannot be the same")
		success = false
	}

	if success {
		var user1 = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Friends[0]).First(&user1)
		var user2 = models.User{}
		userModel2 := this.DB.Where("email = ?", param.Friends[1]).First(&user2)

		if userModel1.RecordNotFound() {
			returnStatus = http.StatusNotFound
			message = append(message, fmt.Sprintf("Email %s not found", param.Friends[0]))
		}

		if userModel2.RecordNotFound() {
			returnStatus = http.StatusNotFound
			message = append(message, fmt.Sprintf("Email %s not found", param.Friends[1]))
			success = false
		}

		if success {
			var userCommon []models.Connection
			connectionModel := this.DB.
				Where("user_id = ? and exists(select 'x' from connections c1 where c1.user_id = ? and c1.friend_id = connections.friend_id)", user1.Id, user2.Id).
					Preload("FriendDetail").Find(&userCommon)

			if connectionModel.RecordNotFound() || len(userCommon) == 0 {
				returnStatus = http.StatusNotFound
				message = append(message, "No common friend found")
				success = false
			} else {
				listEmail := []string{}
				for _, connection := range userCommon {
					listEmail = append(listEmail, connection.FriendDetail.Email)
				}

				ctx.JSON(iris.Map{
					"friends": listEmail,
					"count": len(userCommon),
					"success": success,
				})
				return
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": message,
		"success": success,
	})
}