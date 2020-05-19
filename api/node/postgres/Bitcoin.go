package postgres

import (
	"github.com/dmitriy-vas/node/models"
)

func (p postgres) SearchBitcoinAccount(account string) (out *models.Bitcoin, err error) {
	out = new(models.Bitcoin)
	return out, p.Model(out).
		Where("account = ?", account).
		Select()
}

func (p postgres) SearchBitcoinAddress(address string) (out *models.Bitcoin, err error) {
	out = new(models.Bitcoin)
	return out, p. Model(out).
		Where("address = ?", address).
		Select()
}

func (p postgres) AddBitcoinAccount(account *models.Bitcoin) error {
	_, err := p.Model(account).SelectOrInsert()
	return err
}
