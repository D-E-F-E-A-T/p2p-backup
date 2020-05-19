package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

const (
	SessionName = "_io_auth"
)

var unauthorizedEndpoints = []string{
	"",
	"/",
	"/js/*filepath",
	"/css/*filepath",
	"/img/*filepath",
	"/data/*filepath",
	"/Account/Login",
	"/Account/Logout",
	"/Account/Register",
	"/Account/Restore",
	"/Account/Recovery",
	"/Account/Activate",
	"/Account/Key",
	"/GetOffers",
	"/GetOffersOld",
	"/deal",
	"/deals",
	"/error",
	"/offer",
	"/offers",
	"/verified",
	"/recovery_login",
	"/test",
}

func TokenChecker(c *gin.Context) {
	path := c.FullPath()
	for _, endpoint := range unauthorizedEndpoints {
		if endpoint == path {
			return
		}
	}

	if CheckAuth(c) {
		c.Next()
	}
}

func CheckAuth(c *gin.Context) (authenticated bool) {
	sess := sessions.Default(c)
	token, ok := sess.Get("Token").(string)
	if !ok {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		//c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	session, err := postgres.Postgres.SearchUserSession(token)
	if err == pg.ErrNoRows ||
		(err == nil && session.Expires.Before(time.Now())) {
		c.Redirect(http.StatusTemporaryRedirect, "/")
		//c.AbortWithStatus(http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("Error with checking token: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	ban, err := postgres.Postgres.SearchUserBan(session.ID)
	if err != nil {
		if err != pg.ErrNoRows {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"error": fmt.Sprintf("You have been banned for: %s", ban.Description),
		})
		return
	}
	return true
}

func IsAuthorized(c *gin.Context) {
	c.Status(http.StatusOK)
}
