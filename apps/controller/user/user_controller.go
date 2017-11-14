package user

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
	"social-api/apps/models"
	"time"
	"social-api/apps/helper"
)

type Controller struct {
	DB *gorm.DB
}

// View list user by criteria
func (this *Controller) Index(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages []string
	
	var listUser []models.User

	userParam := ctx.URLParams()

	query := this.DB

	if len(userParam["email"]) > 0 {
		query = query.Where("email like ?", "%"+userParam["email"]+"%")
	}

	if len(userParam["name"]) > 0 {
		query = query.Where("name like ?", "%"+userParam["name"]+"%")
	}

	if len(userParam["status"]) > 0 {
		query = query.Where("status = ?", userParam["status"])
	}
	query.Find(&listUser)

	if len(listUser) > 0 {
		ctx.JSON(iris.Map{
			"users": listUser,
			"success": success,
		})
		return
	} else {
		returnStatus, success = helper.RecordNotFoundMessage("User", &messages)
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": messages,
		"success": success,
	})
}

// View detail user by id
func (this *Controller) Show(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages []string

	id, err := ctx.Params().GetInt("id")

	if err != nil {
		returnStatus, success = helper.InvalidFormatMessage("Id", "Integer", &messages)
	}

	if success {
		var user models.User

		if this.DB.First(&user, id).RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage("User", &messages)
		}

		if success {
			var userOutput models.UserOutput
			userOutput.Id = user.Id
			userOutput.Name = user.Name
			userOutput.Email = user.Email
			if user.Status {
				userOutput.Status = "true"
			} else {
				userOutput.Status = "false"
			}
			userOutput.CreatedAt = user.CreatedAt
			userOutput.UpdatedAt = user.UpdatedAt

			ctx.JSON(iris.Map{
				"user": userOutput,
				"success": success,
			})
			
			return
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": messages,
		"success": success,
	})
}

// Create new user
func (this *Controller) Create(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages []string

	var userParam models.UserOutput
	ctx.ReadJSON(&userParam)

	if len(userParam.Email) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Email", &messages)
	}

	if len(userParam.Name) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Name", &messages)
	}

	if helper.ValidateEMail(userParam.Email) && len(userParam.Email) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Email", "email(mail@domain.com)", &messages)
	}

	if success {
		userModel := this.DB.Where("email = ?", userParam.Email).First(&models.User{})

		if userModel.RecordNotFound() {
			var user models.User
			user.Name = userParam.Name
			user.Email = userParam.Email
			user.Status = true
			user.UpdatedAt = time.Now()
			user.CreatedAt = time.Now()
			
			if err := this.DB.Create(&user).Error; err != nil {
				returnStatus, success = helper.UndefinedErrorMessage(err.Error(), &messages)
			} else {
				ctx.JSON(iris.Map{
					"user":    user,
					"message": "User created successfully",
					"success": true,
				})
				return
			}
		} else {
			returnStatus, success = helper.DuplicateErrorMessage("Email", userParam.Email, &messages)
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": messages,
		"success": success,
	})
}

// update user by user id
func (this *Controller) Update(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages []string

	var userParam models.UserOutput
	ctx.ReadJSON(&userParam)
	
	id, err := ctx.Params().GetInt("id")

	if err != nil {
		returnStatus, success = helper.InvalidFormatMessage("Id", "integer", &messages)
	}

	if success {
		var user = models.User{
			Id: int32(id),
		}

		userModel := this.DB.First(&user, user.Id)

		if userModel.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage("User", &messages)
		} else {
			if len(userParam.Name) > 0 {
				user.Name = userParam.Name
			}

			if len(userParam.Status) > 0 {
				user.Status = userParam.Status == "1"
			}

			if err := this.DB.Save(&user).Error; err != nil {
				returnStatus, success = helper.UndefinedErrorMessage("Failed to update user. " + err.Error(), &messages)
			} else {
				ctx.JSON(iris.Map{
					"message": "User updated successfully",
					"user":    user,
					"success": success,
				})

				return
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": messages,
		"success": success,
	})
}

// Remove user by id
func (this *Controller) Remove(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages []string

	id, err := ctx.Params().GetInt("id")

	if err != nil {
		returnStatus, success = helper.InvalidFormatMessage("Id", "integer", &messages)
	}

	if success {
		var user = models.User{
			Id: int32(id),
		}

		userModel := this.DB.First(&user, user.Id)

		if userModel.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage("User", &messages)
		} else {
			if err := this.DB.Delete(user).Error; err != nil {
				returnStatus, success = helper.UndefinedErrorMessage("Failed to delete user. " + err.Error(), &messages)
			} else {
				ctx.JSON(iris.Map{
					"message": "User deleted successfully",
					"success": true,
				})

				return
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": messages,
		"success": success,
	})
}
