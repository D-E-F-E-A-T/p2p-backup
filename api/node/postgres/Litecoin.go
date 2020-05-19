package postgres

import (
	"github.com/dmitriy-vas/node/models"
)

func (p postgres) SearchLitecoinAccount(account string) (out *models.Litecoin, err error) {
	out = new(models.Litecoin)
	return out, p.Model(out).
		Where("account = ?", account).
		Select()
}

func (p postgres) SearchLitecoinAddress(address string) (out *models.Litecoin, err error) {
	out = new(models.Litecoin)
	return out, p. Model(out).
		Where("address = ?", address).
		Select()
}

func (p postgres) AddLitecoinAccount(account *models.Litecoin) error {
	_, err := p.Model(account).SelectOrInsert()
	return err
}
