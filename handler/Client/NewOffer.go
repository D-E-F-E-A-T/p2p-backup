package Client

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

func NewOffer(c *gin.Context) {
	if !c.GetBool("IsLogged") {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	countries, _ := postgres.Postgres.GetCountries("ru")

	c.HTML(http.StatusOK, "new_offer.html", gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
		"Countries":           countries,
	})
}
