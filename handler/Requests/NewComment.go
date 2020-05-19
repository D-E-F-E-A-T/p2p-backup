package Requests

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type NewCommentRequest struct {
	Deal    uint64 `form:"deal"`
	Message string `form:"message" binding:"required"`
	Rating  uint8  `form:"rating" binding:"min=0,max=10"`
}

func NewComment(c *gin.Context) {
	var request NewCommentRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	deal, err := postgres.Postgres.GetDeal(request.Deal)
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

	var targetUser uint64
	var ownerUser uint64

	if user.ID == deal.Offer.UserID {
		targetUser = deal.Deal.UserID
		ownerUser = deal.Offer.UserID
	} else if user.ID == deal.Deal.UserID {
		targetUser = deal.Offer.UserID
		ownerUser = deal.Deal.UserID
	} else {
		c.Status(http.StatusForbidden)
		return
	}

	if err := postgres.Postgres.AddNewComment(&models.Comment{
		ID:        targetUser,
		UserID:    ownerUser,
		Deal:      deal.ID,
		Message:   request.Message,
		Rating:    request.Rating,
		Timestamp: time.Now(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postgres.Postgres.IncreaseUserRating(targetUser, request.Rating); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusOK)
}
