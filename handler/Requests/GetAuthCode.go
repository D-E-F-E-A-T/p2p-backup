package Requests

import (
	"encoding/hex"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type GetAuthCodeRequest struct {
	Purpose models.AuthPurpose `json:"purpose"`
}

func GetAuthCode(c *gin.Context) {
	var request GetAuthCodeRequest
	if err := c.Bind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	bs := make([]byte, 10)
	rand.Read(bs)

	auth := &models.Auth{
		ID:      user.ID,
		Token:   hex.EncodeToString(bs),
		Purpose: request.Purpose,
		Expires: time.Now().Add(time.Hour * 1),
	}

	if err := postgres.Postgres.AddAuthCode(auth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// TODO send telegram auth code
}
