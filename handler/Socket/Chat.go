package Socket

import (
	"encoding/json"
	"time"

	"gopkg.in/olahol/melody.v1"

	"github.com/dmitriy-vas/p2p/database/postgres"
	"github.com/dmitriy-vas/p2p/handler/middleware"
	"github.com/dmitriy-vas/p2p/models"
)

type Message struct {
	Chat    uint64 `json:"chat"`
	Content string `json:"content"`
}

func (h socketHandler) HandleMessage(session *melody.Session, data []byte) {
	var message Message
	if err := json.Unmarshal(data, &message); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	chat, err := postgres.Postgres.GetChat(message.Chat)
	if err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}
	if chat.Closed {
		return
	}

	user := session.MustGet("User").(*models.User)
	if !middleware.IsAdmin(user.Groups) && user.ID != chat.FirstUser && user.ID != chat.SecondUser {
		return
	}

	sentMessage := &models.Message{
		ChatID:    chat.ID,
		UserID:    user.ID,
		Content:   message.Content,
		Timestamp: time.Now(),
	}

	if err := postgres.Postgres.AddNewMessage(sentMessage); err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	raw, err := json.Marshal(struct {
		Type uint8 `json:"type"`
		*models.Message
		Login string `json:"login"`
	}{
		Type:    TypeMessage,
		Message: sentMessage,
		Login:   user.Login,
	})
	if err != nil {
		session.CloseWithMsg([]byte(err.Error()))
		return
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client, u := range h.Clients {
		if middleware.IsAdmin(u.Groups) || u.ID == chat.FirstUser || u.ID == chat.SecondUser {
			client.Write(raw)
		}
	}
}
