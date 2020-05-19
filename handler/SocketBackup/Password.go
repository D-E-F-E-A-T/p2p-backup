package SocketBackup

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	socketio "github.com/dmitriy-vas/go-socket.io"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/l18n"
	"github.com/dmitriy-vas/p2p/models"
)

func init() {
	Server.OnEvent("", "change_password", OnChangePassword)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"min=8,max=48"`
	NewPassword string `json:"new_password" validate:"min=8,max=48"`
}

func OnChangePassword(s socketio.Conn, c *gin.Context, msg string) {
	var request ChangePasswordRequest
	if err := json.Unmarshal([]byte(msg), &request); err != nil {
		s.Emit("err", err.Error())
		return
	}
	if err := middleware.Valid.Struct(request); err != nil {
		s.Emit("err", err.Error())
		return
	}

	sess := sessions.Default(c)
	userInterface := sess.Get("User")
	if userInterface == nil {
		return
	}

	user := userInterface.(*models.User)

	hash := md5.New()
	hash.Write([]byte(request.OldPassword))
	encodedOldPassword := hex.EncodeToString(hash.Sum(nil))

	if encodedOldPassword != user.Password {
		s.Emit("err", l18n.T("Password is incorrect",
			postgres.LanguageEnglish))
		return
	}

	hash = md5.New()
	hash.Write([]byte(request.NewPassword))
	encodedNewPassword := hex.EncodeToString(hash.Sum(nil))

	if err := postgres.Postgres.SetUserPassword(user.ID, encodedNewPassword); err != nil {
		s.Emit("err", err.Error())
		return
	}
}
