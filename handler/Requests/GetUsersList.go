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

type GetUsersListRequest struct {
	Limit  int    `form:"limit" binding:"min=5,max=50"`
	Page   int    `form:"page" binding:"min=1"`
	Search string `form:"search" binding:"omitempty"`
}

func GetUsersList(c *gin.Context) {
	var request GetUsersListRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)
	if !middleware.IsAdmin(user.Groups) {
		c.Status(http.StatusForbidden)
		return
	}

	count, users, err := postgres.Postgres.GetUsers(request.Search,
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
		"count": count,
		"users": users,
	})
}
