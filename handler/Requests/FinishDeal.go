package Requests

import (
	"encoding/json"
	"fmt"
	"log"
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
)

type FinishDealRequest struct {
	ID uint64 `form:"id" binding:"required"`
}

func FinishDeal(c *gin.Context) {
	var request FinishDealRequest
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

	if (user.ID != deal.Deal.UserID && user.ID != deal.Offer.UserID) ||
		deal.Finished {
		c.Status(http.StatusForbidden)
		return
	}

	var myCurrency uint8
	var isOwner = user.ID == deal.Offer.UserID
	if isOwner {
		myCurrency = deal.Offer.HaveCurrency
	} else {
		myCurrency = deal.Offer.WantCurrency
	}

	if !middleware.IsCryptoUint(myCurrency) {
		c.Status(http.StatusForbidden)
		return
	}

	myCurrencyName := middleware.Currency(myCurrency)
	wrapper := crypto.Wrappers[myCurrencyName]
	if wrapper == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("%s is not available", myCurrencyName),
		})
		return
	}

	var toUserID uint64
	if isOwner {
		toUserID = deal.Deal.UserID
	} else {
		toUserID = deal.Offer.UserID
	}
	login, err := postgres.Postgres.GetUserLogin(toUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postgres.Postgres.SetDealFinished(deal.Deal.ID, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	postgres.Postgres.IncrementFinishedDeals(deal.Deal.UserID)
	postgres.Postgres.IncrementFinishedDeals(deal.Offer.UserID)

	bigAmount := big.NewFloat(float64(deal.Deal.FixedAmount))
	id, err := wrapper.SendTransaction(crypto.SendTransactionRequest{
		From:   user.Login,
		To:     login,
		Amount: bigAmount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Sent transaction, ID: %s", id)

	// TODO add fee
	if err := postgres.Postgres.AddUserTransaction(&models.Transaction{
		ID:        toUserID,
		UserID:    user.ID,
		Deal:      deal.Deal.ID,
		Currency:  myCurrency,
		Amount:    deal.Deal.FixedAmount,
		Timestamp: time.Now(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postgres.Postgres.UnlockUserBalance(user.ID, myCurrency, deal.Deal.FixedAmount); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	msg := &models.Message{
		ChatID:    deal.Chat.ID,
		UserID:    user.ID,
		Content:   "Сделка завершена",
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
		if middleware.IsAdmin(user.Groups) || user.ID == deal.Chat.FirstUser || user.ID == deal.Chat.SecondUser {
			client.Write(raw)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
