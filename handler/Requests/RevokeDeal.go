package Requests

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

type RevokeDealRequest struct {
	Deal uint64 `form:"deal" binding:"required"`
}

func RevokeDeal(c *gin.Context) {
	var request RevokeDealRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)
	if !middleware.IsAdmin(user.Groups) {
		c.Status(http.StatusForbidden)
		return
	}

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
	if deal.Argue == nil {
		c.Status(http.StatusConflict)
		return
	}

	haveCrypto := middleware.IsCryptoUint(deal.Offer.HaveCurrency)
	if haveCrypto && deal.FromApproved {
		err = postgres.Postgres.UnlockUserBalance(deal.Offer.UserID, deal.Offer.HaveCurrency, deal.FixedAmount)
	} else if !haveCrypto && deal.ToApproved {
		err = postgres.Postgres.UnlockUserBalance(deal.Deal.UserID, deal.Offer.WantCurrency, deal.FixedAmount)
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	postgres.Postgres.SetDealFinished(deal.ID, false)
	postgres.Postgres.SetDealApproved(deal.ID, "From", false)
	postgres.Postgres.SetDealApproved(deal.ID, "To", false)
	postgres.Postgres.SetDealAccepted(deal.ID, false)

	if err := postgres.Postgres.DeleteArgue(deal.Argue.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	chats, err := postgres.Postgres.GetArgueChats(deal.Argue.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	for _, chat := range chats {
		if err := postgres.Postgres.CloseChat(chat.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
}
