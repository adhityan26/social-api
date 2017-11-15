// Package block is used to handle
// block user update by email address
package block

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

type blockOutput struct {
	Requestor string `json: requestor`
	Target    string `json: target`
}

// Block user by email address
func (this *Controller) Create(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages = []string{}

	param := blockOutput{}
	ctx.ReadJSON(&param)

	if len(param.Requestor) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Requestor", &messages)
	}

	if len(param.Target) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Target", &messages)
	}

	if helper.ValidateEMail(param.Requestor) && len(param.Requestor) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Requestor", "email(mail@domain.com)", &messages)
	}

	if helper.ValidateEMail(param.Target) && len(param.Target) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Target", "email(mail@domain.com)", &messages)
	}

	if len(param.Requestor) > 0 && len(param.Target) > 0 && param.Requestor == param.Target {
		returnStatus, success = helper.CustomPreconditionErrorMessage("Requestor and Target email cannot be the same", &messages)
	}

	if success {
		var user1 = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Requestor).First(&user1)
		var user2 = models.User{}
		userModel2 := this.DB.Where("email = ?", param.Target).First(&user2)

		if userModel1.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Requestor), &messages)
		}

		if userModel2.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Target), &messages)
		}

		if success {
			var checkBlock models.Block
			checkBlockModel := this.DB.
				Where("requestor_id = ?", user1.Id).
				Where("target_id = ?", user2.Id).
				First(&checkBlock)

			if checkBlockModel.RecordNotFound() {
				if success {
					userBlock := models.Block{}
					userBlock.RequestorId = user1.Id
					userBlock.TargetId = user2.Id
					userBlock.CreatedAt = time.Now()
					userBlock.UpdatedAt = time.Now()
					if err := this.DB.Create(&userBlock).Error; err != nil {
						returnStatus, success = helper.UndefinedErrorMessage(err.Error(), &messages)
					}

					if success {
						ctx.JSON(iris.Map{
							"success": true,
						})
						return
					}
				}
			} else {
				returnStatus, success = helper.CustomPreconditionErrorMessage(fmt.Sprintf("Email %s is already blocked by %s", param.Requestor, param.Target), &messages)
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"messages": messages,
		"success": success,
	})
}

// Block user by email address
func (this *Controller) Remove(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages = []string{}

	param := blockOutput{}
	ctx.ReadJSON(&param)

	if len(param.Requestor) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Requestor", &messages)
	}

	if len(param.Target) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Target", &messages)
	}

	if helper.ValidateEMail(param.Requestor) && len(param.Requestor) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Requestor", "email(mail@domain.com)", &messages)
	}

	if helper.ValidateEMail(param.Target) && len(param.Target) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Target", "email(mail@domain.com)", &messages)
	}

	if len(param.Requestor) > 0 && len(param.Target) > 0 && param.Requestor == param.Target {
		returnStatus, success = helper.CustomPreconditionErrorMessage("Requestor and Target email cannot be the same", &messages)
	}

	if success {
		var user1 = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Requestor).First(&user1)
		var user2 = models.User{}
		userModel2 := this.DB.Where("email = ?", param.Target).First(&user2)

		if userModel1.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Requestor), &messages)
		}

		if userModel2.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Target), &messages)
		}

		if success {
			var checkBlock models.Block
			checkBlockModel := this.DB.
				Where("requestor_id = ?", user1.Id).
				Where("target_id = ?", user2.Id).
				First(&checkBlock)

			if !checkBlockModel.RecordNotFound() {
				if err := checkBlockModel.Delete(models.Block{}).Error; err != nil {
					returnStatus, success = helper.UndefinedErrorMessage(err.Error(), &messages)
				}

				if success {
					ctx.JSON(iris.Map{
						"success": true,
					})
					return
				}
			} else {
				returnStatus, success = helper.CustomPreconditionErrorMessage(fmt.Sprintf("Email %s is not blocked by %s", param.Requestor, param.Target), &messages)
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"messages": messages,
		"success": success,
	})
}
