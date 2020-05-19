package Requests

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/crypto"
	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/Socket"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
	"github.com/dmitriy-vas/p2p/quotas"
)

type ApproveDealRequest struct {
	ID   uint64 `form:"id" binding:"required"`
	Code string `form:"code_google" binding:"omitempty"`
}

func ApproveDeal(c *gin.Context) {
	var request ApproveDealRequest
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

	if !deal.Deal.Accepted || deal.Deal.Finished {
		c.Status(http.StatusForbidden)
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	if user.ID != deal.Deal.UserID && user.ID != deal.Offer.UserID {
		c.Status(http.StatusForbidden)
		return
	}

	haveCurrency := middleware.Currency(deal.Offer.HaveCurrency)
	wantCurrency := middleware.Currency(deal.Offer.WantCurrency)
	var myCurrency uint8
	var isOwner = user.ID == deal.Offer.UserID
	if isOwner {
		myCurrency = deal.Offer.HaveCurrency
	} else {
		myCurrency = deal.Offer.WantCurrency
	}

	var place string
	if myCurrency == deal.Offer.HaveCurrency {
		place = "From"
	} else {
		place = "To"
	}

	// Если моя валюта крипта
	if middleware.IsCryptoUint(myCurrency) {
		//settings, _ := postgres.Postgres.GetUserSettings(user.ID)
		//if settings.Settings.AuthenticatorKey != "" {
		//	if !Account.CheckAuthenticator(request.Code, settings.Settings.AuthenticatorKey) {
		//		c.Status(http.StatusForbidden)
		//		return
		//	}
		//}

		var amountToCheck float32
		if !isOwner {
			amountToCheck = deal.Deal.Amount
		} else {
			cost := quotas.GetCostWithProfit(wantCurrency, haveCurrency, 0)
			amountToCheck = deal.Deal.Amount * cost
		}
		fee, err := postgres.Postgres.GetServiceFee()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		feeAmount := amountToCheck * fee

		myCurrencyName := middleware.Currency(myCurrency)

		wrapper := crypto.Wrappers[myCurrencyName]
		if wrapper == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("%s is not available", myCurrencyName),
			})
			return
		}
		balance, err := wrapper.GetBalance(user.Login)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// TODO add fee
		amountToCheckBig := big.NewFloat(float64(amountToCheck + feeAmount))
		if balance.Cmp(amountToCheckBig) < 0 {
			c.Status(http.StatusConflict)
			return
		}

		// TODO add fee
		if err := postgres.Postgres.SetFixedAndFee(deal.Deal.ID, amountToCheck, feeAmount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		// TODO add fee
		if err := postgres.Postgres.LockUserBalance(user.ID, myCurrency, amountToCheck); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if err := postgres.Postgres.SetDealApproved(deal.Deal.ID, place, true); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		if err := postgres.Postgres.SetDealApproved(deal.Deal.ID, place, true); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	msg := &models.Message{
		UserID:    user.ID,
		ChatID:    deal.Chat.ID,
		Content:   fmt.Sprintf("Подтверждаю перевод %s на счёт", middleware.Currency(myCurrency)),
		Timestamp: time.Now(),
	}
	postgres.Postgres.AddNewMessage(msg)
	raw, _ := json.Marshal(struct {
		Type uint8 `json:"type"`
		*models.Message
		Login string `json:"login"`
	}{
		Type:    Socket.TypeMessage,
		Message: msg,
		Login:   user.Login,
	})

	for client, user := range Socket.SocketHandler.Clients {
		if user.ID == deal.Chat.FirstUser || user.ID == deal.Chat.SecondUser {
			client.Write(raw)
		}
	}
	c.Status(http.StatusOK)
}
