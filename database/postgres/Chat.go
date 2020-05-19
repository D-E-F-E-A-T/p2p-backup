package postgres

import (
	"fmt"

	"github.com/dmitriy-vas/p2p/models"
)

type MessageUser struct {
	tableName struct{} `pg:"users"`
	ID        uint64   `json:"id" pg:",pk"`
	Login     string   `json:"login" pg:",unique,notnull"`
}

type Message struct {
	*models.Message `pg:",inherit"`
	Attachments     []*models.Attachment `json:"attachments"`
	User            *MessageUser         `json:"user"`
}

func (p postgres) GetChatMessages(id uint64, limit, offset int) (count int, messages []*Message, err error) {
	count, err = p.Model(&messages).
		Relation("User").
		Relation("Attachments").
		Where("chat_id = ?", id).
		Order("timestamp DESC").
		Limit(limit).
		Offset(offset).
		SelectAndCount()
	return
}

func (p postgres) AddMessageAttachment(attachment *models.Attachment) error {
	return p.Insert(attachment)
}

func (p postgres) AddNewMessage(message *models.Message) error {
	return p.Insert(message)
}

func (p postgres) DeleteMessages(id uint64) error {
	var m []*models.Message
	_, err := p.Model(&m).
		Where("message.deal = ?", id).
		Delete()
	return err
}

type GetUserChatsModel struct {
	*models.Chat `pg:",inherit"`
	FromUser     struct {
		ID    uint64 `json:"id"`
		Login string `json:"login"`
	} `json:"first_user"`
	ToUser struct {
		ID    uint64 `json:"id"`
		Login string `json:"login"`
	} `json:"second_user"`
	Argue *models.Argue `json:"argue"`
}

// TODO get relation through builder
// TODO delete checks for admin
func (p postgres) GetUserChats(limit, offset int, sortMethod string, category uint8, user *models.User) (count int, chats []*GetUserChatsModel, err error) {
	selectText := `"chat"."id",
			       "chat"."deal_id",
			       "chat"."first_user",
			       "chat"."second_user",
			       "chat"."first_unread",
			       "chat"."second_unread",
			       "chat"."argue_id",
			       "chat"."closed",
			       "chat"."timestamp",
                   "first_user"."id" AS "from_user__id",
			       "first_user"."login" AS "from_user__login",
                   "second_user"."id" AS "to_user__id",
			       "second_user"."login" AS "to_user__login",
			       "argue"."id"       AS "argue__id",
			       "argue"."deal_id"  AS "argue__deal_id",
			       "argue"."category" AS "argue__category",
			       "argue"."finished" AS "argue__finished"`

	joinText := `FROM "chats" AS "chat"
				 LEFT JOIN "argues" AS "argue" ON "argue"."id" = "chat"."argue_id"
		         LEFT JOIN "users" AS "first_user" ON "first_user"."id" = "chat"."first_user"
		         LEFT JOIN "users" AS "second_user" ON "second_user"."id" = "chat"."second_user"`
	switch sortMethod {
	case "New":
		sortMethod = "DESC"
	case "Old":
		sortMethod = "ASC"
	}
	endText := fmt.Sprintf(`ORDER BY "timestamp" %s
				LIMIT %d
				OFFSET %d;`, sortMethod, limit, offset)

	queryText := "WHERE (chat.closed IS FALSE)"
	if category != 0 {
		queryText += fmt.Sprintf("\nAND (argue.category = %d)", category)
	} else {
		queryText += fmt.Sprintf("\nAND (argue_id IS NULL)")
	}
	queryText += fmt.Sprintf("\nAND ((chat.first_user = %d) OR (chat.second_user = %d))", user.ID, user.ID)

	_, err = p.Query(&chats, fmt.Sprintf("SELECT %s\n %s\n %s\n %s", selectText, joinText, queryText, endText))
	if err != nil {
		return count, chats, err
	}

	_, err = p.Query(&count, fmt.Sprintf("SELECT count(*)\n %s\n %s", joinText, queryText))
	return
}

func (p postgres) AddUserChat(chat *models.Chat) error {
	return p.Insert(chat)
}

func (p postgres) GetChat(id uint64) (chat *models.Chat, err error) {
	chat = new(models.Chat)
	return chat, p.Model(chat).
		Where("id = ?", id).
		Select()
}

type GetChatWithDealModel struct {
	*models.Chat `pg:",inherit"`
	Deal         *models.Deal `json:"deal"`
}

func (p postgres) GetChatWithDeal(id uint64, argue bool, user uint64) (result *GetChatWithDealModel, err error) {
	result = new(GetChatWithDealModel)
	query := p.Model(result).
		Relation("Deal").
		Where("deal.id = ?", id)
	if argue {
		query = query.Where("chat.argue_id IS NOT NULL").
			Where("chat.first_user = 0").Where("chat.second_user = ?", user)
	} else {
		query.Where("chat.argue IS NULL").
			Where("chat.first_user != 0").Where("chat.second_user != 0")
	}
	return result, query.Select()
}

func (p postgres) CloseChat(id uint64) error {
	_, err := p.Model((*models.Chat)(nil)).
		Where("id = ?", id).
		Set("closed = ?", true).
		Update()
	return err
}

func (p postgres) GetArgueChats(id uint64) (chats []*models.Chat, err error) {
	return chats, p.Model(&chats).
		Where("argue_id = ?", id).
		Select()
}

type GetArguesChatsModel struct {
	*models.Argue `pg:",inherit"`
	FromChat      *models.Chat `json:"first_chat"`
	ToChat        *models.Chat `json:"second_chat"`
	FromUser      *models.User `json:"first_user"`
	ToUser        *models.User `json:"second_user"`
}

func (p postgres) GetArguesChats(limit, offset int, category uint8, sortMethod string) (count int, argues []*GetArguesChatsModel, err error) {
	selectText := `argue.*,
       from_chat.id           as from_chat__id,
       from_chat.second_user  as from_chat__second_user,
       from_chat.first_unread as from_chat__first_unread,
       from_chat.timestamp    as from_chat__timestamp,
       to_chat.id             as to_chat__id,
       to_chat.second_user    as to_chat__second_user,
       to_chat.first_unread   as to_chat__first_unread,
       to_chat.timestamp      as to_chat__timestamp,
       from_user.login        as from_user__login,
	   from_user.id           as from_user__id,
       to_user.login          as to_user__login,
	   to_user.id             as to_user__id`
	joinText := `FROM "argues" "argue"
         LEFT JOIN chats "from_chat"
                   on from_chat.second_user = argue.first_user AND from_chat.argue_id = argue.id
         LEFT JOIN chats "to_chat"
                   on to_chat.second_user = argue.second_user AND to_chat.argue_id = argue.id
         LEFT JOIN users "from_user"
                   on from_user.id = argue.first_user
         LEFT JOIN users "to_user"
                   on to_user.id = argue.second_user`
	queryText := fmt.Sprintf("WHERE (argue.finished IS FALSE) AND (argue.category = %d)", category)
	switch sortMethod {
	case "New":
		sortMethod = "DESC"
	case "Old":
		sortMethod = "ASC"
	}
	endText := fmt.Sprintf("ORDER BY argue.last_activity %s LIMIT %d OFFSET %d;", sortMethod, limit, offset)
	_, err = p.Query(&argues, fmt.Sprintf("SELECT %s\n %s\n %s\n %s", selectText, joinText, queryText, endText))
	if err != nil {
		return count, argues, err
	}
	_, err = p.Query(&count, fmt.Sprintf("SELECT COUNT(*) %s\n %s", joinText, queryText))
	return
}
