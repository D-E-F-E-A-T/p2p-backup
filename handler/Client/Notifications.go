package Client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Notifications(c *gin.Context) {
	c.HTML(http.StatusOK, "notifications.html", gin.H{
		"NumOffers": c.GetInt("NumOffers"),
		"NumDeals":  c.GetInt("NumDeals"),
		"IsLogged":  c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
	})
}
