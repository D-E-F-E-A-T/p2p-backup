package Account

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/spf13/viper"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
	"github.com/dmitriy-vas/p2p/smtp"
)

type RestoreRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}

func Restore(c *gin.Context) {
	var request RestoreRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	restore, err := postgres.Postgres.SearchRestoreToken(request.Token)
	if err != nil {
		if err == pg.ErrNoRows {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	bs := make([]byte, 14)
	rand.Read(bs)
	password := hex.EncodeToString(bs)

	hash := md5.New()
	hash.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hash.Sum(nil))

	email, err := postgres.Postgres.SearchUserAndReturnEmail(restore.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postgres.Postgres.RestoreUserPassword(restore.ID, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postgres.Postgres.DeleteRestoreToken(restore.Token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/recovery_login")
	go func() {
		var buf bytes.Buffer
		smtp.Mail.Template.Execute(&buf, smtp.MailOption{
			PreText:     fmt.Sprintf("New password: %s", password),
			AfterText:   "",
			ButtonLink:  "",
			ButtonText:  "",
			Company:     viper.GetString("smtp.company"),
			Unsubscribe: "",
		})

		smtp.Mail.Send(email,
			`New password created`,
			buf.String())
	}()
}

type RecoveryRequest struct {
	Email string `json:"email" form:"email" binding:"email"`
}

func Recovery(c *gin.Context) {
	var request RecoveryRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := postgres.Postgres.SearchUser(request.Email)
	if err != nil {
		if err == pg.ErrNoRows {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	bs := make([]byte, 24)
	rand.Read(bs)
	token := hex.EncodeToString(bs)

	restore := &models.Restore{
		ID:      user.ID,
		Token:   token,
		Expires: time.Now().AddDate(0, 0, 1),
	}

	if err := postgres.Postgres.AddNewRestore(restore); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
	go func() {
		var buf bytes.Buffer
		smtp.Mail.Template.Execute(&buf, smtp.MailOption{
			PreText:   "",
			AfterText: "",
			ButtonLink: fmt.Sprintf("%s/Account/Restore?token=%s",
				viper.GetString("api.host"),
				token),
			ButtonText:  "Restore access to account",
			Company:     viper.GetString("smtp.company"),
			Unsubscribe: "",
		})

		smtp.Mail.Send(user.Email,
			`Restore access`,
			buf.String())
	}()
}
