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

type EditOfferRequest struct {
	OfferID      uint64         `form:"offer_id" binding:"required"`
	WantCurrency uint8          `form:"want_currency" binding:"required,nefield=HaveCurrency"`
	HaveCurrency uint8          `form:"have_currency" binding:"required,nefield=WantCurrency"`
	Location     string         `form:"location" binding:"omitempty,len=2"`
	Provider     uint8          `form:"provider" binding:"required"`
	Profit       int16          `form:"profit" binding:"required_without=Cost"`
	Cost         float32        `form:"cost" binding:"required_without=Profit"`
	Minimal      float32        `form:"minimal" binding:"min=0,ltfield=Maximal"`
	Maximal      float32        `form:"maximal" binding:"min=0,gtfield=Minimal"`
	Title        string         `form:"title" binding:"omitempty,max=20"`
	Terms        string         `form:"terms" binding:"omitempty,max=256"`
	Days         []time.Weekday `form:"days[]" binding:"dive,min=0,max=6"`
	StartTime    time.Time      `form:"start_time" binding:"required_with=FinishTime,ltfield=FinishTime" time_format:"15:04"`
	FinishTime   time.Time      `form:"finish_time" binding:"required_with=StartTime,gtfield=StartTime" time_format:"15:04"`
	Active       bool           `form:"active" binding:"omitempty"`
	Warranty     bool           `form:"warranty" binding:"omitempty"`
	WithPhone    bool           `form:"with_phone" binding:"omitempty"`
	WithDocs     bool           `form:"with_docs" binding:"omitempty"`
	WithDeals    uint           `form:"with_deals" binding:"omitempty"`
	WithRating   uint           `form:"with_rating" binding:"omitempty"`
	CancelTime   uint8          `form:"cancel_time"`
}

func EditOffer(c *gin.Context) {
	var request EditOfferRequest
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

	offer.WantCurrency = request.WantCurrency
	offer.HaveCurrency = request.HaveCurrency
	offer.Location = request.Location
	offer.Provider = request.Provider
	offer.Profit = request.Profit
	offer.Cost = request.Cost
	offer.Minimal = request.Minimal
	offer.Maximal = request.Maximal
	offer.Title = request.Title
	offer.Terms = request.Terms
	offer.Days = request.Days
	offer.StartTime = request.StartTime
	offer.StartTime = time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		request.StartTime.Hour(),
		request.StartTime.Minute(),
		0,
		0,
		time.Now().Location())
	offer.FinishTime = time.Date(time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		request.FinishTime.Hour(),
		request.FinishTime.Minute(),
		0,
		0,
		time.Now().Location())
	offer.Warranty = request.Warranty
	offer.WithPhone = request.WithPhone
	offer.WithDocs = request.WithDocs
	offer.WithDeals = request.WithDeals
	offer.WithRating = request.WithRating
	offer.CancelTime = request.CancelTime

	if err := postgres.Postgres.UpdateOffer(offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusOK)
}
