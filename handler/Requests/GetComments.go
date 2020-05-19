package Requests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

type GetCommentsRequest struct {
	User       uint64 `form:"user" binding:"required"`
	SortMethod string `form:"sort_method" binding:"oneof=New Old"`
	Limit      int    `form:"limit" binding:"min=5,max=50"`
	Page       int    `form:"page" binding:"min=1"`
}

func GetComments(c *gin.Context) {
	var request GetCommentsRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	count, comments, err := postgres.Postgres.SearchUserComments(request.User,
		request.Limit,
		(request.Page-1)*request.Limit,
		request.SortMethod)
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
		"total":    count,
		"comments": comments,
	})
}
