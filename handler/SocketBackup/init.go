package SocketBackup

import (
	"log"
	"time"

	socketio "github.com/dmitriy-vas/go-socket.io"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

var Server, _ = socketio.NewServer(nil)

func init() {
	Server.OnConnect("", OnConnect)
	Server.OnDisconnect("", OnDisconnect)
	Server.OnError("", OnError)
}

func OnConnect(s socketio.Conn, c *gin.Context) error {
	sess := sessions.Default(c)
	if _, ok := sess.Get("Token").(string); !ok {
		return s.Close()
	}

	user, ok := sess.Get("User").(*models.User)
	if !ok {
		return s.Close()
	}

	postgres.Postgres.UpdateLastSeen(user.ID, time.Now())
	return nil
}

func OnDisconnect(s socketio.Conn, c *gin.Context, i string) {
	sess := sessions.Default(c)
	user := sess.Get("User").(*models.User)

	postgres.Postgres.UpdateLastSeen(user.ID, time.Now())
}

func OnError(s socketio.Conn, c *gin.Context, err error) {
	if err.Error() == "websocket: close 1001 (going away)" {
		return
	}

	log.Printf("Client encountered error: %v", err)
}
