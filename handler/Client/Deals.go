package Client

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/models"
)

func Deals(c *gin.Context) {
	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	c.HTML(http.StatusOK, "deals.html", gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
		"User":                user,
	})
}
