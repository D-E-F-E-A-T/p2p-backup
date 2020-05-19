package Socket

import (
	"encoding/json"

	"gopkg.in/olahol/melody.v1"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

type NewFeeRequest struct {
	Fee float32 `json:"fee"`
}

func (h socketHandler) HandleNewFee(session *melody.Session, data []byte) {
	var request NewFeeRequest
	if err := json.Unmarshal(data, &request); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	user := session.MustGet("User").(*models.User)
	if !middleware.IsAdmin(user.Groups) {
		return
	}

	if err := postgres.Postgres.SetServiceFee(request.Fee); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}
}
