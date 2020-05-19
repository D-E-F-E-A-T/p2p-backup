package Requests

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type DeleteOfferRequest struct {
	OfferID uint64 `form:"offer_id" binding:"required"`
}

func DeleteOffer(c *gin.Context) {
	var request DeleteOfferRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	offer, err := postgres.Postgres.GetOffer(request.OfferID)
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

	if offer.UserID != user.ID {
		c.Status(http.StatusForbidden)
		return
	}

	count, err := postgres.Postgres.SearchOfferDeals(offer.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	if count > 0 {
		c.Status(http.StatusConflict)
		return
	}

	if err := postgres.Postgres.DeleteOffer(offer.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}
