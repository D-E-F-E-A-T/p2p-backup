package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) LockUserBalance(user uint64, currency uint8, balance float32) error {
	wallet := &models.Wallet{
		UserID:   user,
		Currency: currency,
		Balance:  balance,
	}
	_, err := p.Model(wallet).
		Where("user_id = ?", wallet.UserID).
		Where("currency = ?", wallet.Currency).
		OnConflict("(user_id, currency) DO UPDATE").
		Set("balance = COALESCE(wallet.balance,0) + ?", wallet.Balance).
		Insert()
	return err
}

func (p postgres) UnlockUserBalance(user uint64, currency uint8, balance float32) error {
	_, err := p.Model((*models.Wallet)(nil)).
		Where("user_id = ?", user).
		Where("currency = ?", currency).
		Set("balance = wallet.balance - ?", balance).
		Update()
	return err
}
