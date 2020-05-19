package Requests

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/quotas"
)

type GetUserOffersRequest struct {
	ID         uint64 `form:"id" binding:"required"`
	SortMethod string `form:"sort_method" binding:"oneof=Popularity New Old"`
	Limit      int    `form:"limit" binding:"min=5,max=50"`
	Page       int    `form:"page" binding:"min=1"`
}

func GetUserOffers(c *gin.Context) {
	var request GetUserOffersRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	count, offers, err := postgres.Postgres.SearchOffers("",
		0,
		0,
		"",
		"",
		0,
		request.SortMethod,
		false,
		false,
		0,
		request.ID,
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
			haveCurrency := middleware.Currency(offer.HaveCurrency)
			wantCurrency := middleware.Currency(offer.WantCurrency)
			offer.Cost = quotas.Quotas[haveCurrency][wantCurrency]
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"count":  count,
		"offers": offers,
	})
}
