package Client

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/models"
)

func Offers(c *gin.Context) {
	if !c.GetBool("IsLogged") {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	c.HTML(http.StatusOK, "offers.html", gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
		"User":                user,
	})
}
