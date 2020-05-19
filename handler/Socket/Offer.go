package Socket

import (
	"encoding/json"

	"gopkg.in/olahol/melody.v1"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/models"
)

type Toggle struct {
	ID uint64 `json:"id"`
}

func (h socketHandler) HandleToggle(session *melody.Session, data []byte) {
	var toggle Toggle
	if err := json.Unmarshal(data, &toggle); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	offer, err := postgres.Postgres.GetOffer(toggle.ID)
	if err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	user := session.MustGet("User").(*models.User)
	if user.ID != offer.UserID {
		session.Close()
		return
	}

	if err := postgres.Postgres.ToggleOffer(offer.ID, offer.UserID); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
	}
}
