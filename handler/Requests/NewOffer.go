package Requests

import (
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type NewOfferRequest struct {
	WantCurrency uint8          `form:"want_currency" binding:"required,nefield=HaveCurrency,currency"`
	HaveCurrency uint8          `form:"have_currency" binding:"required,nefield=WantCurrency,currency"`
	Location     string         `form:"location" binding:"country"`
	Provider     uint8          `form:"provider" binding:"provider"`
	Profit       int16          `form:"profit" binding:"min=-99"`
	Cost         float32        `form:"cost"`
	Minimal      float32        `form:"minimal" binding:"omitempty,min=0,required_with=Maximal,ltfield=Maximal"`
	Maximal      float32        `form:"maximal" binding:"omitempty,min=0,required_with=Minimal,gtfield=Minimal"`
	Title        string         `form:"title" binding:"omitempty,max=40"`
	Terms        string         `form:"terms" binding:"omitempty,max=256"`
	Days         []time.Weekday `form:"days[]" binding:"required,dive,min=0,max=6"`
	StartTime    time.Time      `form:"start_time" binding:"required_with=FinishTime,ltfield=FinishTime" time_format:"15:04"`
	FinishTime   time.Time      `form:"finish_time" binding:"required_with=StartTime,gtfield=StartTime" time_format:"15:04"`
	Active       bool           `form:"active" binding:"omitempty"`
	Warranty     bool           `form:"warranty" binding:"omitempty"`
	WithPhone    bool           `form:"with_phone" binding:"omitempty"`
	WithDocs     bool           `form:"with_docs" binding:"omitempty"`
	WithDeals    uint           `form:"with_deals" binding:"omitempty"`
	WithRating   uint           `form:"with_rating" binding:"omitempty"`
	CancelTime   uint8          `form:"cancel_time" binding:"min=60"`
}

func NewOffer(c *gin.Context) {
	var request NewOfferRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	offer := &models.Offer{
		UserID:       user.ID,
		StatisticID:  user.ID,
		SettingsID:   user.ID,
		WantCurrency: request.WantCurrency,
		HaveCurrency: request.HaveCurrency,
		Minimal:      request.Minimal,
		Maximal:      request.Maximal,
		Profit:       request.Profit,
		Cost:         request.Cost,
		Title:        request.Title,
		Location:     request.Location,
		Provider:     request.Provider,
		Terms:        request.Terms,
		Warranty:     request.Warranty,
		WithDynamic:  request.Cost == 0,
		WithPhone:    request.WithPhone,
		WithDocs:     request.WithDocs,
		WithDeals:    request.WithDeals,
		WithRating:   request.WithRating,
		Active:       request.Active,
		CancelTime:   request.CancelTime,
		Days:         request.Days,
		StartTime: time.Date(time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			request.StartTime.Hour(),
			request.StartTime.Minute(),
			0,
			0,
			time.Now().Location()),
		FinishTime: time.Date(time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			request.FinishTime.Hour(),
			request.FinishTime.Minute(),
			0,
			0,
			time.Now().Location()),
		Timestamp: time.Now(),
	}

	if err := postgres.Postgres.AddNewOffer(offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": offer.ID,
	})
}
