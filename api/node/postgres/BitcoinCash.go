package postgres

import (
	"github.com/dmitriy-vas/node/models"
)

func (p postgres) SearchBitcoinCashAccount(account string) (out *models.BitcoinCash, err error) {
	out = new(models.BitcoinCash)
	return out, p.Model(out).
		Where("account = ?", account).
		Select()
}

func (p postgres) SearchBitcoinCashAddress(address string) (out *models.BitcoinCash, err error) {
	out = new(models.BitcoinCash)
	return out, p. Model(out).
		Where("address = ?", address).
		Select()
}

func (p postgres) AddBitcoinCashAccount(account *models.BitcoinCash) error {
	_, err := p.Model(account).SelectOrInsert()
	return err
}
