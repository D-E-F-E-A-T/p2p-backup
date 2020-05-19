package Requests

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/Socket"
	"github.com/dmitriy-vas/p2p/models"
)

type AcceptDealRequest struct {
	ID uint64 `form:"id" binding:"required"`
}

func AcceptDeal(c *gin.Context) {
	var request AcceptDealRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	deal, err := postgres.Postgres.GetDeal(request.ID)
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

	if user.ID != deal.Offer.UserID || deal.Accepted {
		c.Status(http.StatusForbidden)
		return
	}

	if err := postgres.Postgres.SetDealAccepted(request.ID, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	msg := &models.Message{
		UserID:    user.ID,
		ChatID:    deal.Chat.ID,
		Content:   "Сделка активирована",
		Timestamp: time.Now(),
	}
	postgres.Postgres.AddNewMessage(msg)
	raw, _ := json.Marshal(struct {
		Type uint8 `json:"type"`
		*models.Message
		Login string `json:"login"`
	}{
		Type:    Socket.TypeMessage,
		Message: msg,
		Login:   user.Login,
	})

	for client, user := range Socket.SocketHandler.Clients {
		if user.ID == deal.Chat.FirstUser || user.ID == deal.Chat.SecondUser {
			client.Write(raw)
		}
	}
}
