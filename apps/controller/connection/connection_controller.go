// Package connection is used to handle
// user connection friend by email address
package connection

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
	"fmt"
	"social-api/apps/models"
	"time"
	"social-api/apps/helper"
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

// View list friends by email address
func (this *Controller) Index(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages = []string{}

	param := friendList{}
	ctx.ReadJSON(&param)

	if len(param.Email) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Email", &messages)
	}

	if helper.ValidateEMail(param.Email) && len(param.Email) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Email", "email(mail@domain.com)", &messages)
	}

	if success {
		var user = models.User{}
		userModel := this.DB.Where("email = ?", param.Email).
			Preload("Friends").
			Preload("Friends.UserDetail").First(&user)
		if userModel.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage("User", &messages)
		} else {
			listEmail := []string{}
			for _, friend := range user.Friends {
				listEmail = append(listEmail, friend.UserDetail.Email)
			}

			ctx.JSON(iris.Map{
				"friends": listEmail,
				"count":   len(user.Friends),
				"success": success,
			})
			return
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"messages": messages,
		"success":  success,
	})
}

// Create connection between two email address
func (this *Controller) Create(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages = []string{}

	param := connectionOutput{}
	ctx.ReadJSON(&param)

	if len(param.Friends) < 2 {
		returnStatus, success = helper.CustomPreconditionRequiredErrorMessage("Must provide 2 email", &messages)
	}

	if (len(param.Friends) == 2) && (param.Friends[0] == param.Friends[1]) {
		returnStatus, success = helper.CustomPreconditionErrorMessage("Email cannot be the same", &messages)
	}

	if len(param.Friends) > 0 && helper.ValidateEMail(param.Friends[0]) && len(param.Friends[0]) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Email 1", "email(mail@domain.com)", &messages)
	}

	if len(param.Friends) > 1 && helper.ValidateEMail(param.Friends[1]) && len(param.Friends[1]) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Email 2", "email(mail@domain.com)", &messages)
	}

	if success {
		var user1 = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Friends[0]).First(&user1)
		var user2 = models.User{}
		userModel2 := this.DB.Where("email = ?", param.Friends[1]).First(&user2)

		if userModel1.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Friends[0]), &messages)
		}

		if userModel2.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Friends[1]), &messages)
		}

		if success {
			var checkBlocked1 models.Block
			checkBlockedModel1 := this.DB.
				Where("requestor_id = ? and target_id = ?", user1.Id, user2.Id).
				First(&checkBlocked1)

			var checkBlocked2 models.Block
			checkBlockedModel2 := this.DB.
				Where("requestor_id = ? and target_id = ?", user2.Id, user1.Id).
				First(&checkBlocked2)

			if checkBlockedModel1.RecordNotFound() && checkBlockedModel2.RecordNotFound() {
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
					tx := this.DB.Begin()
					userConnection1 := models.Connection{}
					userConnection1.UserId = user1.Id
					userConnection1.FriendId = user2.Id
					userConnection1.CreatedAt = time.Now()
					userConnection1.UpdatedAt = time.Now()
					if err := tx.Create(&userConnection1).Error; err != nil {
						returnStatus, success = helper.UndefinedErrorMessage(err.Error(), &messages)
					}

					if success {
						userConnection2 := models.Connection{}
						userConnection2.UserId = user2.Id
						userConnection2.FriendId = user1.Id
						userConnection2.CreatedAt = time.Now()
						userConnection2.UpdatedAt = time.Now()
						if err := tx.Create(&userConnection2).Error; err != nil {
							returnStatus, success = helper.UndefinedErrorMessage(err.Error(), &messages)
						}
					}

					if success {
						tx.Commit()
						ctx.JSON(iris.Map{
							"success": true,
						})
						return
					} else {
						tx.Rollback()
					}
				} else {
					returnStatus, success = helper.CustomPreconditionErrorMessage(fmt.Sprintf("Email %s and %s is already friend", param.Friends[0], param.Friends[1]), &messages)
				}
			} else {
				if !checkBlockedModel1.RecordNotFound() {
					returnStatus, success = helper.CustomPreconditionErrorMessage(fmt.Sprintf("Email %s is blocked by %s", param.Friends[0], param.Friends[1]), &messages)
				}
				if !checkBlockedModel2.RecordNotFound() {
					returnStatus, success = helper.CustomPreconditionErrorMessage(fmt.Sprintf("Email %s is blocked by %s", param.Friends[1], param.Friends[0]), &messages)
				}
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"messages": messages,
		"success":  success,
	})
}

// view common friend between teo emil address
func (this *Controller) Common(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages = []string{}

	param := connectionOutput{}
	ctx.ReadJSON(&param)

	if len(param.Friends) < 2 {
		returnStatus, success = helper.CustomPreconditionErrorMessage("Must provide 2 email", &messages)
	}

	if (len(param.Friends) == 2) && (param.Friends[0] == param.Friends[1]) {
		returnStatus, success = helper.CustomPreconditionErrorMessage("Email cannot be the same", &messages)
	}

	if success {
		var user1 = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Friends[0]).First(&user1)
		var user2 = models.User{}
		userModel2 := this.DB.Where("email = ?", param.Friends[1]).First(&user2)

		if userModel1.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Friends[0]), &messages)
		}

		if userModel2.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Friends[1]), &messages)
		}

		if success {
			var userCommon []models.Connection
			connectionModel := this.DB.
				Where("user_id = ? and exists(select 'x' from connections c1 where c1.user_id = ? and c1.friend_id = connections.friend_id)", user1.Id, user2.Id).
				Preload("FriendDetail").Find(&userCommon)

			if connectionModel.RecordNotFound() || len(userCommon) == 0 {
				returnStatus, success = helper.RecordNotFoundMessage("Common friend", &messages)
			} else {
				listEmail := []string{}
				for _, connection := range userCommon {
					listEmail = append(listEmail, connection.FriendDetail.Email)
				}

				ctx.JSON(iris.Map{
					"friends": listEmail,
					"count":   len(userCommon),
					"success": success,
				})
				return
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"messages": messages,
		"success":  success,
	})
}
