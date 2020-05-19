package SocketBackup

import (
	"strconv"

	socketio "github.com/dmitriy-vas/go-socket.io"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

func init() {
	Server.OnEvent("", "accept_deal", OnAcceptDeal)
}

func OnAcceptDeal(s socketio.Conn, c *gin.Context, i interface{}) {
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

	deal, err := postgres.Postgres.GetDeal(i.(uint64))
	if err != nil {
		s.Emit("err", err.Error())
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)
	_ = user

	if middleware.IsCryptoUint(deal.Offer.HaveCurrency) {

	} else {

	}
}
