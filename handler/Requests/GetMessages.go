package Requests

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

type GetMessagesRequest struct {
	Chat  uint64 `json:"chat" form:"chat" binding:"required"`
	Limit int    `json:"limit" form:"limit" binding:"min=5,max=50"`
	Page  int    `json:"page" form:"page" binding:"min=1"`
}

func GetMessages(c *gin.Context) {
	var request GetMessagesRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	chat, err := postgres.Postgres.GetChat(request.Chat)
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

	if !middleware.IsAdmin(user.Groups) && user.ID != chat.FirstUser && user.ID != chat.SecondUser {
		c.Status(http.StatusForbidden)
		return
	}

	count, messages, err := postgres.Postgres.GetChatMessages(chat.ID,
		request.Limit,
		(request.Page-1)*request.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total":    count,
		"messages": messages,
	})
}
