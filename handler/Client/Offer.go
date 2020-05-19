package Client

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
	"github.com/dmitriy-vas/p2p/quotas"
)

type OfferRequest struct {
	ID     uint64 `form:"id" binding:"required"`
	Status string `form:"status" binding:"omitempty,oneof=Sell Buy"`
}

func Offer(c *gin.Context) {
	if !c.GetBool("IsLogged") {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	var request OfferRequest
	if err := c.Bind(&request); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	offer, err := postgres.Postgres.SearchOffer(request.ID)
	if err != nil {
		if err == pg.ErrNoRows {
			c.Redirect(http.StatusTemporaryRedirect, "/")
		} else {
			c.HTML(http.StatusOK, "error.html", gin.H{
				"error": err.Error(),
			})
		}
		return
	}
	countries, _ := postgres.Postgres.GetCountries("ru")

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	var isOwner = user.ID == offer.UserID
	haveCurrency := middleware.Currency(offer.HaveCurrency)
	wantCurrency := middleware.Currency(offer.WantCurrency)

	if offer.WithDynamic {
		quota := quotas.GetCostWithProfit(haveCurrency, wantCurrency, offer.Profit)
		offer.Cost = quota
	}

	c.HTML(http.StatusOK, "offer.html", gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
		"Offer":               offer,
		"Status":              request.Status,
		"Countries":           countries,
		"IsOwner":             isOwner,
		"Quotas":              quotas.Quotas,
	})

}
