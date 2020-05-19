package Client

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type UserRequest struct {
	ID uint64 `form:"id" binding:"required"`
}

func User(c *gin.Context) {
	if !c.GetBool("IsLogged") {
		c.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	var request UserRequest
	if err := c.Bind(&request); err != nil {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	info, err := postgres.Postgres.GetUserInfo(request.ID)
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
	info.User.ID = request.ID

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	var relationship = &models.Relationship{}
	if user.ID != info.User.ID {
		relationship, err = postgres.Postgres.GetUserRelationship(user.ID, info.User.ID)
		if err != nil && err != pg.ErrNoRows {
			c.HTML(http.StatusOK, "error.html", gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	countries, _ := postgres.Postgres.GetCountries("ru")
	c.HTML(http.StatusOK, "user.html", gin.H{
		"NumOffers":           c.GetInt("NumOffers"),
		"NumDeals":            c.GetInt("NumDeals"),
		"IsLogged":            c.GetBool("IsLogged"),
		"NotificationsAmount": c.GetInt("NotificationsAmount"),
		"Notifications":       c.MustGet("Notifications"),
		"Countries":           countries,
		"User":                info,
		"Relationship":        relationship,
	})
}
