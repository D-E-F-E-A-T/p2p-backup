package SocketBackup

import (
	"log"
	"strconv"

	socketio "github.com/dmitriy-vas/go-socket.io"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

func init() {
	Server.OnEvent("", "toggle_offer", OnToggleOffer)
}

func OnToggleOffer(s socketio.Conn, c *gin.Context, i interface{}) {
	log.Printf("ToggleOfferType: %T", i)

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user, ok := userInterface.(*models.User)
	if !ok {
		return
	}

	switch i.(type) {
	case string:
		id, err := strconv.ParseUint(i.(string), 10, 64)
		if err != nil {
			s.Emit("err", err.Error())
			return
		}
		i = id
	case float64:
		i = uint64(i.(float64))
	}

	if err := postgres.Postgres.ToggleOffer(i.(uint64), user.ID); err != nil {
		s.Emit("err", err.Error())
	}
}
