package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) AddNewUserSession(session *models.Session) error {
	return p.Insert(session)
}

func (p postgres) SearchUserSession(token string) (session *models.Session, err error) {
	session = new(models.Session)
	return session, p.Model(session).
		Where("token = ?", token).
		Select()
}
