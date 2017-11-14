// Package message is used to handle
// sending user message
package message

import (
	"github.com/kataras/iris"
	"github.com/jinzhu/gorm"
	"fmt"
	"social-api/apps/models"
	"time"
	"regexp"
	"social-api/apps/helper"
)

type Controller struct {
	DB *gorm.DB
}

type messageOutput struct {
	Sender string `json: sender`
	Text   string `json: text`
}

// Create message and view list user that can receive update
func (this *Controller) Create(ctx iris.Context) {
	var returnStatus, success = 200, true
	var messages = []string{}

	param := messageOutput{}
	ctx.ReadJSON(&param)

	if len(param.Sender) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Sender", &messages)
	}

	if len(param.Text) == 0 {
		returnStatus, success = helper.MandatoryErrorMessage("Text", &messages)
	}

	if helper.ValidateEMail(param.Sender) && len(param.Sender) > 0 {
		returnStatus, success = helper.InvalidFormatMessage("Sender", "email(mail@domain.com)", &messages)
	}

	if success {
		var user = models.User{}
		userModel1 := this.DB.Where("email = ?", param.Sender).First(&user)

		if userModel1.RecordNotFound() {
			returnStatus, success = helper.RecordNotFoundMessage(fmt.Sprintf("Email %s", param.Sender), &messages)
		}

		if success {
			listMentionMatch := regexp.MustCompile(`([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})`)
			listMentionTmp := listMentionMatch.FindAllString(param.Text, -1)
			var listMention []string
			var listEmailSent []string

			for _, mention := range listMentionTmp {
				if mention != user.Email {
					listMention = append(listMention, mention)
				}
			}

			var followers []models.User
			this.DB.
				Where("(exists(select 'x' from connections c where c.user_id = ? and c.friend_id = users.id) or "+
				"exists(select 'x' from subscribes s where s.target_id = ? and s.requestor_id = users.id) or "+
				"(users.email in (?)))", user.Id, user.Id, listMention).
				Where("not exists(select 'x' from blocks b where b.requestor_id = ? and b.target_id = users.id)", user.Id).
				Find(&followers)

			for _, follower := range followers {
				listEmailSent = append(listEmailSent, follower.Email)
			}

			if len(listEmailSent) > 0 {
				userMessage := models.Message{}
				userMessage.SenderId = user.Id
				userMessage.Text = param.Text
				userMessage.CreatedAt = time.Now()
				userMessage.UpdatedAt = time.Now()
				if err := this.DB.Create(&userMessage).Error; err != nil {
					returnStatus, success = helper.UndefinedErrorMessage(err.Error(), &messages)
				}

				if success {
					ctx.JSON(iris.Map{
						"success":    true,
						"recipients": listEmailSent,
					})
					return
				}
			} else {
				returnStatus, success = helper.RecordNotFoundMessage("Follower/mention", &messages)
			}
		}
	}

	ctx.StatusCode(returnStatus)
	ctx.JSON(iris.Map{
		"message": messages,
		"success": success,
	})
}
