package postgres

import (
	"github.com/dmitriy-vas/node/models"
)

func (p postgres) SearchEthereumAccount(account string) (out *models.Ethereum, err error) {
	out = new(models.Ethereum)
	return out, p.Model(out).
		Where("account = ?", account).
		Select()
}

func (p postgres) AddEthereumAccount(account *models.Ethereum) error {
	return p.Insert(account)
}
