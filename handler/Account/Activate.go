package Account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/database/postgres"
)

type ActivateRequest struct {
	Token string `json:"token" form:"token" binding:"required"`
}

func Activate(c *gin.Context) {
	var request ActivateRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := postgres.Postgres.SearchUserByActivation(request.Token)
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

	if err := postgres.Postgres.UpdateEmailVerified(user.ID, true); err != nil {
		if err == pg.ErrNoRows {
			c.Status(http.StatusNotFound)
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	if err := postgres.Postgres.DeleteUserActivation(user.ID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, "/verified")
}
