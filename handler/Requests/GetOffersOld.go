package Requests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

type GetOffersOldRequest struct {
	HaveCurrency string `json:"have_currency" form:"have_currency" binding:"omitempty,oneof=btc bch eth ltc pzm xrp wvs eur gbp rub usd,nefield=WantCurrency"`
	WantCurrency string `json:"want_currency" form:"want_currency" binding:"omitempty,oneof=btc bch eth ltc pzm xrp wvs eur gbp rub usd,nefield=HaveCurrency"`
	ProviderName string `json:"provider_name" form:"provider_name" binding:"omitempty,oneof=TestProvider"`
	Limit        int    `json:"limit" form:"limit" binding:"min=5,max=50"`
	Page         int    `json:"page" form:"page" binding:"min=1"`
}

func GetOffersOld(c *gin.Context) {
	var request GetOffersOldRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	count, offers, err := postgres.Postgres.SearchOffersOld(request.HaveCurrency,
		request.WantCurrency,
		request.ProviderName,
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

	c.JSON(http.StatusOK, gin.H{
		"count":  count,
		"offers": offers,
	})
}
