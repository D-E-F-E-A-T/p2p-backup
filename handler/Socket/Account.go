package Socket

import (
	"encoding/json"
	"time"

	"gopkg.in/olahol/melody.v1"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

type NewBanRequest struct {
	User        uint64    `json:"user"`
	Expires     time.Time `json:"expires"`
	Description string    `json:"description"`
}

func (h socketHandler) HandleNewBan(session *melody.Session, data []byte) {
	var request NewBanRequest
	if err := json.Unmarshal(data, &request); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	user := session.MustGet("User").(*models.User)
	if !middleware.IsAdmin(user.Groups) {
		return
	}

	if err := postgres.Postgres.AddUserBan(&models.Ban{
		UserID:      request.User,
		Expires:     request.Expires,
		Description: request.Description,
	}); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}
	if err := postgres.Postgres.DeleteSessions(request.User); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}
}

type ChangeGroupsRequest struct {
	User   uint64 `json:"user_id"`
	Groups []int  `json:"groups"`
}

func (h socketHandler) HandleChangeGroups(session *melody.Session, data []byte) {
	var request ChangeGroupsRequest
	if err := json.Unmarshal(data, &request); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	user := session.MustGet("User").(*models.User)
	if !middleware.IsAdmin(user.Groups) {
		return
	}

	if err := postgres.Postgres.UpdateUserGroups(request.User, request.Groups); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}
}
