package Socket

import (
	"encoding/json"

	"gopkg.in/olahol/melody.v1"

	"github.com/dmitriy-vas/p2p/crypto"
	"github.com/dmitriy-vas/p2p/models"
)

type CreateWallet struct {
	Currency string `json:"currency" binding:"crypto"`
}

func (h socketHandler) HandleCreateWallet(session *melody.Session, data []byte) {
	var request CreateWallet
	if err := json.Unmarshal(data, &request); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	user := session.MustGet("User").(*models.User)
	wrapper, ok := crypto.Wrappers[request.Currency]
	if !ok {
		return
	}

	address, err := wrapper.CreateAccount(user.Login)
	if err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}
	raw, err := json.Marshal(struct {
		Type     uint8  `json:"type"`
		Address  string `json:"address"`
		Currency string `json:"currency"`
	}{
		Type:     TypeCreateWallet,
		Address:  address,
		Currency: request.Currency,
	})

	session.Write(raw)
}
