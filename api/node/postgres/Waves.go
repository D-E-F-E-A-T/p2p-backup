package postgres

import (
	"github.com/dmitriy-vas/node/models"
)

func (p postgres) SearchWavesAccount(account string) (out *models.Waves, err error) {
	out = new(models.Waves)
	return out, p.Model(out).
		Where("account = ?", account).
		Select()
}

func (p postgres) AddWavesAccount(account *models.Waves) error {
	return p.Insert(account)
}
