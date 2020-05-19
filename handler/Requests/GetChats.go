package Requests

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

type GetChatsRequest struct {
	Limit      int    `json:"limit" form:"limit" binding:"min=5,max=50"`
	Page       int    `json:"page" form:"page" binding:"min=1"`
	SortMethod string `json:"sort_method" form:"sort_method" binding:"oneof=New Old"`
	Category   uint8  `json:"category" form:"category"`
}

func GetChats(c *gin.Context) {
	var request GetChatsRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	if middleware.IsAdmin(user.Groups) {
		count, argues, err := postgres.Postgres.GetArguesChats(
			request.Limit,
			(request.Page-1)*request.Limit,
			request.Category,
			request.SortMethod)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"count":  count,
			"argues": argues,
		})
	} else {
		count, chats, err := postgres.Postgres.GetUserChats(request.Limit,
			(request.Page-1)*request.Limit,
			request.SortMethod,
			request.Category,
			user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"count": count,
			"chats": chats,
		})
	}
}
