// Package subscribe is used to handle
// user subscription friend by email address
package subscribe

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

type subscribeOutput struct {
	Requestor string `json: requestor`
	Target string `json: target`
}

// create subscription user to receive user updates
func (this *Controller) Create(ctx iris.Context) {
	param := subscribeOutput{}
	ctx.ReadJSON(&param)

	var success, returnStatus = true, 200
	var message = []string{}

	if len(param.Requestor) == 0 {
		returnStatus = http.StatusPreconditionRequired
		message = append(message, "Requestor cannot be empty")
		success = false
	}

	if len(param.Target) == 0 {
		returnStatus = http.StatusPreconditionRequired
		message = append(message, "Target cannot be empty")
		success = false
	}

	if len(param.Requestor) > 0 && len(param.Target) > 0 && param.Requestor == param.Target {
		returnStatus = http.StatusPreconditionFailed
		message = append(message, "Requestor and Target email cannot be the same")
		success = false
	}

	if success {
		var user1 = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Requestor).First(&user1)
		var user2 = models.User{}
		userModel2 := this.DB.Where("email = ?", param.Target).First(&user2)

		if userModel1.RecordNotFound() {
			returnStatus = http.StatusNotFound
			message = append(message, fmt.Sprintf("Email %s not found", param.Requestor))
			success = false
		}

		if userModel2.RecordNotFound() {
			returnStatus = http.StatusNotFound
			message = append(message, fmt.Sprintf("Email %s not found", param.Target))
			success = false
		}

		if success {
			var checkBlocked1 models.Block
			checkBlockedModel1 := this.DB.
				Where("requestor_id = ? and target_id = ?", user2.Id, user1.Id).
				First(&checkBlocked1)

			if checkBlockedModel1.RecordNotFound() {

				var checkSubscribe models.Subscribe
				checkSubscribeModel := this.DB.
					Where("requestor_id = ?", user1.Id).
					Where("target_id = ?", user2.Id).
					First(&checkSubscribe)

				if checkSubscribeModel.RecordNotFound() {
					if success {
						userSubscribe := models.Subscribe{}
						userSubscribe.RequestorId = user1.Id
						userSubscribe.TargetId = user2.Id
						userSubscribe.CreatedAt = time.Now()
						userSubscribe.UpdatedAt = time.Now()
						if err := this.DB.Create(&userSubscribe).Error; err != nil {
							returnStatus = http.StatusInternalServerError
							message = append(message, err.Error())
							success = false
						}

						if success {
							ctx.JSON(iris.Map{
								"success": true,
							})
							return
						}
					}
				} else {
					returnStatus = http.StatusPreconditionFailed
					message = append(message, fmt.Sprintf("Email %s already subscribe %s", param.Requestor, param.Target))
					success = false
				}
			} else {
				if !checkBlockedModel1.RecordNotFound() {
					returnStatus = http.StatusPreconditionFailed
					message = append(message, fmt.Sprintf("Email %s is blocked by %s", param.Requestor, param.Target))
					success = false
				}
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": message,
		"success": success,
	})
}