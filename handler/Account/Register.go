package Account

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
	"github.com/dmitriy-vas/p2p/smtp"
)

type RegisterRequest struct {
	Login    string `json:"login" form:"login" binding:"min=5,max=32"`
	Email    string `json:"email" form:"email" binding:"email"`
	Phone    string `json:"phone" form:"phone"`
	Password string `json:"password" form:"password" binding:"min=8,max=48"`
}

func Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	hash := md5.New()
	hash.Write([]byte(request.Password))
	encodedPassword := hex.EncodeToString(hash.Sum(nil))

	bs := make([]byte, 64)
	rand.Read(bs)
	activation := hex.EncodeToString(bs)

	user := &models.User{
		Login:      request.Login,
		Email:      request.Email,
		Phone:      request.Phone,
		Password:   encodedPassword,
		Activation: activation,
	}

	if err := postgres.Postgres.RegisterNewUser(user); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
	go func() {
		var buf bytes.Buffer
		smtp.Mail.Template.Execute(&buf, smtp.MailOption{
			PreText:   "Account activation",
			AfterText: "",
			ButtonLink: fmt.Sprintf("%s/Account/Activate?token=%s",
				viper.GetString("api.host"),
				activation),
			ButtonText:  "Activate",
			Company:     viper.GetString("smtp.company"),
			Unsubscribe: "",
		})

		smtp.Mail.Send(user.Email,
			`Account activation`,
			buf.String(),
		)
	}()
}
