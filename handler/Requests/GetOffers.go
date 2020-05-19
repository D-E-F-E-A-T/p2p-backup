package Requests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/quotas"
)

type GetOffersRequest struct {
	OfferMethod    string  `json:"offer_method" form:"offer_method" binding:"omitempty,oneof=SellOrBuy Sell Buy"`
	CryptoCurrency uint8   `json:"crypto_currency" form:"crypto_currency" binding:"omitempty,crypto"`
	Currency       uint8   `json:"currency" form:"currency" binding:"omitempty,fiat"`
	Location       string  `json:"location" form:"location" binding:"omitempty,country"`
	Search         string  `json:"search" form:"search" binding:"omitempty,max=40"`
	Amount         float32 `json:"amount" form:"amount" binding:"omitempty,min=0"`
	Warranty       bool    `json:"warranty" form:"warranty" binding:"omitempty"`
	Provider       uint8   `json:"provider" form:"provider" binding:"omitempty,provider"`
	SortMethod     string  `json:"sort_method" form:"sort_method" binding:"oneof=New Old"`
	Limit          int     `json:"limit" form:"limit" binding:"min=5,max=50"`
	Page           int     `json:"page" form:"page" binding:"min=1"`
}

func GetOffers(c *gin.Context) {
	var request GetOffersRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	count, offers, err := postgres.Postgres.SearchOffers(request.OfferMethod,
		request.CryptoCurrency,
		request.Currency,
		request.Location,
		request.Search,
		request.Amount,
		request.SortMethod,
		request.Warranty,
		true,
		request.Provider,
		0,
		request.Limit,
		(request.Page-1)*request.Limit)
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

	for _, offer := range offers {
		if offer.WithDynamic {
			offer.Cost = quotas.GetCostWithProfit(offer.HaveCurrency, offer.WantCurrency, offer.Profit)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count":  count,
		"offers": offers,
	})
}
