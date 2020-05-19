package SocketBackup

import (
	socketio "github.com/dmitriy-vas/go-socket.io"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

func init() {
	Server.OnEvent("", "add_available_address", AddAvailableAddress)
	Server.OnEvent("", "del_available_address", DelAvailableAddress)
}

func AddAvailableAddress(s socketio.Conn, c *gin.Context, msg string) {
	sess := sessions.Default(c)
	if _, ok := sess.Get("Token").(string); !ok {
		return
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.Var(msg, "ipv4"); err != nil {
			s.Emit("err", err.Error())
			return
		}
	}

	userInterface := sess.Get("User")
	user := userInterface.(*models.User)
	addresses, err := postgres.Postgres.GetUserAvailableAddresses(user.ID)
	if err != nil {
		s.Emit("err", err.Error())
		return
	}

	for _, adr := range addresses {
		if adr == msg {
			return
		}
	}
	addresses = append(addresses, msg)

	if err := postgres.Postgres.SetUserAvailableAddresses(user.ID, addresses); err != nil {
		s.Emit("err", err.Error())
	}
}

func DelAvailableAddress(s socketio.Conn, c *gin.Context, msg string) {
	sess := sessions.Default(c)
	if _, ok := sess.Get("Token").(string); !ok {
		return
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.Var(msg, "ipv4"); err != nil {
			s.Emit("err", err.Error())
			return
		}
	}

	userInterface := sess.Get("User")
	user := userInterface.(*models.User)
	addresses, err := postgres.Postgres.GetUserAvailableAddresses(user.ID)
	if err != nil {
		s.Emit("err", err.Error())
		return
	}

	for i, adr := range addresses {
		if adr == msg {
			addresses = append(addresses[:i], addresses[i:]...)
			break
		}
	}

	if err := postgres.Postgres.SetUserAvailableAddresses(user.ID, addresses); err != nil {
		s.Emit("err", err.Error())
	}
}
