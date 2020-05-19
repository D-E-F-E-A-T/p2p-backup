package Client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

type EditOfferRequest struct {
	ID uint64 `form:"id" binding:"required"`
}

func EditOffer(c *gin.Context) {
	if !c.GetBool("IsLogged") {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	var request EditOfferRequest
	if err := c.Bind(&request); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	offer, err := postgres.Postgres.GetOffer(request.ID)
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

	c.HTML(http.StatusOK, "edit_offer.html", gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
		"Offer":               offer,
		"Countries":           countries,
	})
}
