// Package block is used to handle
// block user update by email address
package block

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

type blockOutput struct {
	Requestor string `json: requestor`
	Target string `json: target`
}

// Block user by email address
func (this *Controller) Create(ctx iris.Context) {
	param := blockOutput{}
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
				message = append(message, fmt.Sprintf("Email %s already blocked %s", param.Requestor, param.Target))
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