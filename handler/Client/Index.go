package Client

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

func Index(c *gin.Context) {
	countries, _ := postgres.Postgres.GetCountries("ru")

	c.HTML(http.StatusOK, "page.html", gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
		"Countries":           countries,
	})
}
