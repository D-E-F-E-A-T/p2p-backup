package Requests

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	"github.com/spf13/viper"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/Socket"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

func UploadFile(c *gin.Context) {
	chatStr, ok := c.GetPostForm("chat")
	if !ok {
		c.Status(http.StatusBadRequest)
		return
	}
	log.Printf("ChatSTR: %s", chatStr)
	chatID, err := strconv.ParseUint(chatStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	d, err := postgres.Postgres.GetChat(chatID)
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

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	if !middleware.IsAdmin(user.Groups) && user.ID != d.FirstUser && user.ID != d.SecondUser {
		c.Status(http.StatusForbidden)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		return
	}

	msg := &models.Message{
		ChatID:    chatID,
		UserID:    user.ID,
		Content:   " ",
		Timestamp: time.Now(),
	}
	postgres.Postgres.AddNewMessage(msg)

	var savedFiles []*models.Attachment
	for _, file := range files {
		filename := filepath.Base(file.Filename)
		path := path.Join(
			viper.GetString("api.files"),
			filename,
		)
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		attachment := &models.Attachment{
			MessageID: msg.ID,
			Path:      filepath.Join(viper.GetString("api.host"), "/files/", filename),
			Name:      filename,
		}
		postgres.Postgres.AddMessageAttachment(attachment)
		savedFiles = append(savedFiles, attachment)
	}

	raw, _ := json.Marshal(struct {
		Type uint8 `json:"type"`
		*models.Message
		Attachments []*models.Attachment `json:"attachments"`
		Login       string               `json:"login"`
	}{
		Type:        Socket.TypeMessage,
		Message:     msg,
		Attachments: savedFiles,
		Login:       user.Login,
	})
	for client, user := range Socket.SocketHandler.Clients {
		if middleware.IsAdmin(user.Groups) || user.ID == d.FirstUser || user.ID == d.SecondUser {
			client.Write(raw)
		}
	}
}
