package Requests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

type GetUserDealsRequest struct {
	ID         uint64 `form:"id" binding:"required"`
	SortMethod string `form:"sort_method" binding:"oneof=New Old"`
	Limit      int    `form:"limit" binding:"min=5,max=50"`
	Page       int    `form:"page" binding:"min=1"`
}

func GetUserDeals(c *gin.Context) {
	var request GetUserDealsRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	count, deals, err := postgres.Postgres.GetUserDeals(request.ID,
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
		"count": count,
		"deals": deals,
	})
}
