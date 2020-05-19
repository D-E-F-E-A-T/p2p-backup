package Requests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/Socket"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

// TODO add category and message validation
type NewArgueRequest struct {
	ID       uint64 `form:"id" binding:"required"`
	Message  string `form:"message" binding:"required"`
	Category uint8  `form:"category" binding:"required"`
}

func NewArgue(c *gin.Context) {
	var request NewArgueRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err := postgres.Postgres.GetDeal(request.ID)
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
	if data.Argue != nil {
		c.Status(http.StatusConflict)
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	if (user.ID != data.Deal.UserID && user.ID != data.Offer.UserID) ||
		(data.Finished) ||
		(!data.Accepted) {
		c.Status(http.StatusForbidden)
		return
	}

	go func() {
		notification := &models.Notification{
			Title:     "New argue created",
			Link:      "http://95.169.186.48:9999/chat",
			Content:   fmt.Sprintf("%s: new argue for deal %d", user.Login, data.ID),
			Timestamp: time.Now(),
		}
		notification.ID = data.UserID
		postgres.Postgres.AddNewNotification(notification)
		notification.ID = data.Offer.UserID
		postgres.Postgres.AddNewNotification(notification)
		for client, user := range Socket.SocketHandler.Clients {
			if user.ID == data.UserID || user.ID == data.Offer.UserID {
				notification.ID = user.ID
				raw, _ := json.Marshal(struct {
					Type uint8 `json:"type"`
					*models.Notification
				}{
					Type:         Socket.TypeNotification,
					Notification: notification,
				})
				client.Write(raw)
			}
		}
	}()

	argue := &models.Argue{
		DealID:       request.ID,
		LastActivity: time.Now(),
		FirstUser:    data.Deal.UserID,
		SecondUser:   data.Offer.UserID,
		Category:     request.Category,
	}
	if err := postgres.Postgres.AddNewArgue(argue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg := &models.Message{
		UserID:    user.ID,
		Content:   request.Message,
		Timestamp: time.Now(),
	}
	sendMsg := func(chat uint64) {
		msg.ChatID = chat
		postgres.Postgres.AddNewMessage(msg)
	}

	chat := &models.Chat{
		DealID:     data.Deal.ID,
		SecondUser: data.Deal.UserID,
		ArgueID:    argue.ID,
		Timestamp:  time.Now(),
	}

	postgres.Postgres.AddUserChat(chat)
	if user.ID == chat.SecondUser {
		sendMsg(chat.ID)
	}

	chat.ID = 0
	chat.SecondUser = data.Offer.UserID
	postgres.Postgres.AddUserChat(chat)
	if user.ID == chat.SecondUser {
		sendMsg(chat.ID)
	}

	raw, _ := json.Marshal(struct {
		Type uint8 `json:"type"`
		*models.Argue
		//FirstChat  *models.Chat `json:"first_chat"`
		//SecondChat *models.Chat `json:"second_chat"`
		//FirstUser  *models.User `json:"first_user"`
		//SecondUser *models.User `json:"second_user"`
	}{
		Type:  Socket.TypeNewArgue,
		Argue: argue,
	})
	for client, user := range Socket.SocketHandler.Clients {
		if middleware.IsAdmin(user.Groups) || user.ID == data.Deal.UserID || user.ID == data.Offer.UserID {
			client.Write(raw)
		}
	}

	c.Status(http.StatusOK)
}
