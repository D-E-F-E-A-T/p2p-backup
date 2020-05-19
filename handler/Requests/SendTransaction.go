package Requests

import (
	"math/big"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/crypto"
	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

type SendTransactionRequest struct {
	Destination string  `form:"destination"`
	Currency    string  `form:"currency" binding:"crypto"`
	Amount      float32 `form:"amount"`
}

func SendTransaction(c *gin.Context) {
	var request SendTransactionRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	wrapper := crypto.Wrappers[request.Currency]
	if wrapper == nil {
		c.Status(http.StatusForbidden)
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	balance, err := wrapper.GetBalance(user.Login)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": err.Error(),
		})
		return
	}

	bigAmount := big.NewFloat(float64(request.Amount))
	if balance.Cmp(bigAmount) < 0 {
		c.Status(http.StatusConflict)
		return
	}

	id, err := wrapper.SendTransaction(crypto.SendTransactionRequest{
		From:   user.Login,
		To:     request.Destination,
		Amount: bigAmount,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := postgres.Postgres.AddUserTransaction(&models.Transaction{
		UserID:    user.ID,
		Currency:  middleware.CurrencyID(request.Currency),
		Amount:    request.Amount,
		Timestamp: time.Now(),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": id,
	})
}
