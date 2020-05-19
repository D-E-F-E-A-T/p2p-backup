package Requests

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type GetNotificationsRequest struct {
	Limit int `json:"limit" form:"limit" binding:"min=5,max=50"`
	Page  int `json:"page" form:"page" binding:"min=1"`
}

func GetNotifications(c *gin.Context) {
	var request GetNotificationsRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user, _ := userInterface.(*models.User)

	count, notifications, err := postgres.Postgres.SearchUserNotifications(user.ID,
		request.Limit,
		(request.Page-1)*request.Limit)
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

	c.JSON(http.StatusOK, gin.H{
		"count":         count,
		"notifications": notifications,
	})
}
