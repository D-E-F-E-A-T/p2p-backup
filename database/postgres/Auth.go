package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) AddAuthCode(auth *models.Auth) error {
	return p.Insert(auth)
}
