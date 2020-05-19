package Socket

import (
	"encoding/json"
	"sync"

	"gopkg.in/olahol/melody.v1"

	"github.com/dmitriy-vas/p2p/models"
)

type socketHandler struct {
	mutex   *sync.Mutex
	Clients map[*melody.Session]*models.User
}

var SocketHandler *socketHandler

func init() {
	SocketHandler = &socketHandler{
		mutex:   new(sync.Mutex),
		Clients: make(map[*melody.Session]*models.User),
	}
}

func (h socketHandler) Connect(session *melody.Session) {
	SocketHandler.mutex.Lock()
	defer SocketHandler.mutex.Unlock()
	user := session.MustGet("User").(*models.User)
	SocketHandler.Clients[session] = user
}

func (h socketHandler) Disconnect(session *melody.Session) {
	SocketHandler.mutex.Lock()
	defer SocketHandler.mutex.Unlock()
	delete(SocketHandler.Clients, session)
}

type PreMessage struct {
	Type uint8 `json:"type"`
}

const (
	TypeMessage = iota
	TypeToggle
	TypeCreateWallet
	TypeNotification
	TypeNewArgue
	TypeNewFee
	TypeNewBan
	TypeChangeGroups
)

func (h socketHandler) Handle(mrouter *melody.Melody, session *melody.Session, data []byte) {
	var preMessage PreMessage
	if err := json.Unmarshal(data, &preMessage); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	switch preMessage.Type {
	case TypeMessage:
		h.HandleMessage(session, data)
	case TypeToggle:
		h.HandleToggle(session, data)
	case TypeCreateWallet:
		h.HandleCreateWallet(session, data)
	case TypeNewFee:
		h.HandleNewFee(session, data)
	case TypeNewBan:
		h.HandleNewBan(session, data)
	case TypeChangeGroups:
		h.HandleChangeGroups(session, data)
	}
}
