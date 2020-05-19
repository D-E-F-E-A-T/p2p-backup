package SocketBackup

import (
	socketio "github.com/dmitriy-vas/go-socket.io"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/handler/middleware"
)

func init() {
	Server.OnEvent("", "get_balance", OnGetBalance)
}

func OnGetBalance(s socketio.Conn, c *gin.Context, msg string) {
	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	if userInterface == nil {
		return
	}

	if err := middleware.Valid.Var(msg, "oneof=btc bch eth ltc pzm xrp wvs"); err != nil {
		s.Emit("err", err.Error())
		return
	}

	// TODO finish getting balance from nodes
}
