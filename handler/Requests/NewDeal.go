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

type NewDealRequest struct {
	Offer   uint64  `form:"offer" binding:"required"`
	Amount  float32 `form:"amount" binding:"required,min=0"`
	Message string  `form:"message" binding:"required"`
}

func NewDeal(c *gin.Context) {
	var request NewDealRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	offer, err := postgres.Postgres.GetOffer(request.Offer)
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
	userInterface := sess.Get("User").(*models.User)

	if userInterface.ID == offer.UserID || !offer.Active {
		c.Status(http.StatusForbidden)
		return
	}

	user, err := postgres.Postgres.GetUserInfo(userInterface.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	user.ID = userInterface.ID

	now := time.Now()
	// TODO uncomment on prod
	//isAvailableDay := func() bool {
	//	day := now.Weekday()
	//	for _, d := range offer.Days {
	//		if d == day {
	//			return true
	//		}
	//	}
	//	return false
	//}
	//
	//generateTime := func(t time.Time) time.Time {
	//	return time.Date(now.Year(), now.Month(), now.Day(),
	//		t.Hour(),
	//		t.Minute(),
	//		0,
	//		0,
	//		t.Location())
	//}

	//if (offer.WithDocs && !user.Statistic.DocsVerified) ||
	//	(offer.WithPhone && !user.Statistic.PhoneVerified) ||
	//	(offer.WithRating != 0 && user.Statistic.Rating < offer.WithRating) ||
	//	(offer.WithDeals != 0 && user.Statistic.FinishedDeals < offer.WithDeals) ||
	//	(offer.Minimal != 0 && request.Amount < offer.Minimal) ||
	//	(offer.Maximal != 0 && request.Amount > offer.Maximal) { //||
	//	//(!isAvailableDay()) ||
	//	//(generateTime(offer.StartTime).Before(now)) ||
	//	//(offer.FinishTime.Hour() != 0 && offer.FinishTime.Minute() != 0 && generateTime(offer.FinishTime).After(now)) {
	//	c.Status(http.StatusForbidden)
	//	return
	//}

	// TODO set minimal amount eq or more than 5 USD
	deal := &models.Deal{
		OfferID:   offer.ID,
		UserID:    user.ID,
		Amount:    request.Amount,
		Timestamp: now,
	}

	if err := postgres.Postgres.AddNewDeal(deal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	chat := &models.Chat{
		DealID:     deal.ID,
		FirstUser:  deal.UserID,
		SecondUser: offer.UserID,
		Timestamp:  time.Now(),
	}
	if err := postgres.Postgres.AddUserChat(chat); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postgres.Postgres.AddNewMessage(&models.Message{
		ChatID:    chat.ID,
		UserID:    user.ID,
		Content:   request.Message,
		Timestamp: time.Now(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": deal.ID,
	})
}
