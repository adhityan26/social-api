package user

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
	"social-api/apps/models"
	"time"
	"net/http"
	"fmt"
)

type Controller struct {
	DB *gorm.DB
}

func (this *Controller) Index(ctx iris.Context) {
	var listUser []models.User

	userParam := ctx.URLParams()

	query := this.DB

	if len(userParam["email"]) > 0 {
		query = query.Where("email like ?", "%" + userParam["email"] + "%")
	}

	if len(userParam["name"]) > 0 {
		query = query.Where("name like ?", "%" + userParam["name"] + "%")
	}

	if len(userParam["status"]) > 0 {
		query = query.Where("status = ?", userParam["status"])
	}
	query.Find(&listUser)

	if len(listUser) == 0 {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(iris.Map{
			"message": "User not found",
		})
		return
	}

	ctx.JSON(iris.Map{
		"data": listUser,
	})
}

func (this *Controller) Show(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")

	if err != nil {
		ctx.StatusCode(http.StatusPreconditionFailed)
		ctx.JSON(iris.Map{
			"message":"Invalid format",
			"trace":err.Error(),
		})
		return
	}

	var user models.User

	if this.DB.First(&user, id).RecordNotFound() {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(map[string]string {
			"message":"User not found",
		})
		return
	}

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

	ctx.JSON(userOutput)
}

func (this *Controller) Create(ctx iris.Context) {
	var user models.User
	var userParam models.UserOutput
	ctx.ReadJSON(&userParam)
	user.Name = userParam.Name
	user.Email = userParam.Email
	user.Status = true
	user.UpdatedAt = time.Now()
	user.CreatedAt = time.Now()
	var success, returnStatus = true, 200
	var message = []string{}

	if len(user.Email) == 0 {
		returnStatus = http.StatusPreconditionRequired
		message = append(message, "Email cannot be empty")
		success = false
	}

	if len(user.Name) == 0 {
		returnStatus = http.StatusPreconditionRequired
		message = append(message, "Name cannot be empty")
		success = false
	}

	if success {
		var checkUser = models.User{
			Email: user.Email,
		}

		userModel := this.DB.Where("email = ?", user.Email).First(&checkUser)

		if userModel.RecordNotFound() {
			if err := this.DB.Create(&user).Error; err != nil {
				returnStatus = http.StatusInternalServerError
				message = append(message, err.Error())
				success = false
			} else {
				ctx.JSON(iris.Map{
					"user": user,
					"message": "User created successfully",
					"success": true,
				})
				return
			}
		} else {
			returnStatus = http.StatusConflict
			message = append(message, fmt.Sprintf("Email %s is already exists", user.Email))
			success = false
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": message,
		"success": success,
	})
}

func (this *Controller) Update(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")

	if err != nil {
		ctx.StatusCode(http.StatusPreconditionFailed)
		ctx.JSON(iris.Map{
			"message":"Invalid format",
			"trace":err.Error(),
			"success": false,
		})
		return
	}

	var user = models.User{
		Id: int32(id),
	}

	userModel := this.DB.First(&user, user.Id)

	if userModel.RecordNotFound() {
		ctx.JSON(iris.Map{
			"message": "User not found",
			"success": false,
		})
	} else {
		var userParam models.UserOutput
		ctx.ReadJSON(&userParam)
		if len(userParam.Name) > 0 {
			user.Name = userParam.Name
		}

		if len(userParam.Status) > 0 {
			user.Status = userParam.Status == "1"
		}

		if err := this.DB.Save(&user).Error; err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"message": "Failed to update user",
				"trace": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(iris.Map{
			"message": "User updated successfully",
			"user": user,
			"success": true,
		})
	}
}

func (this *Controller) Remove(ctx iris.Context) {
	id, err := ctx.Params().GetInt("id")

	if err != nil {
		ctx.StatusCode(http.StatusPreconditionFailed)
		ctx.JSON(iris.Map{
			"message":"Invalid format",
			"trace":err.Error(),
			"success": false,
		})
		return
	}

	var user = models.User{
		Id: int32(id),
	}

	userModel := this.DB.First(&user, user.Id)

	if userModel.RecordNotFound() {
		ctx.JSON(iris.Map{
			"message": "User not found",
			"success": false,
		})
		return
	} else {
		if err := this.DB.Delete(user).Error; err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(iris.Map{
				"message": "Failed to delete user",
				"trace": err.Error(),
				"success": false,
			})
			return
		}

		ctx.JSON(iris.Map{
			"message": "User deleted successfully",
			"success": true,
		})
	}
}