package Requests

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type CancelDealRequest struct {
	ID uint64 `form:"id" binding:"required"`
}

func CancelDeal(c *gin.Context) {
	var request CancelDealRequest
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

	// TODO delete messages
	if err := postgres.Postgres.DeleteDeal(request.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.Status(http.StatusOK)
}
