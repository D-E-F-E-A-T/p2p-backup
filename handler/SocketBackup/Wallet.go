package SocketBackup

import (
	socketio "github.com/dmitriy-vas/go-socket.io"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/crypto"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

func init() {
	Server.OnEvent("", "create_wallet", OnCreateWallet)
}

type NewWalletResponse struct {
	Address  string `json:"address"`
	Currency string `json:"currency"`
}

func OnCreateWallet(s socketio.Conn, c *gin.Context, i interface{}) {
	switch i.(type) {
	case float64:
		i = middleware.Currency(uint8(i.(float64)))
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	user := userInterface.(*models.User)

	if wrapper, ok := crypto.Wrappers[i.(string)]; ok {
		address, err := wrapper.CreateAccount(user.Login)
		if err != nil {
			s.Emit("err", err.Error())
			return
		}
		s.Emit("new_wallet", NewWalletResponse{
			Address:  address,
			Currency: i.(string),
		})
	}
}
