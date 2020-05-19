package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) SetServiceFee(fee float32) error {
	_, err := p.Model((*models.ServiceSettings)(nil)).
		Set("fee = ?", fee).
		Update()
	return err
}

func (p postgres) GetServiceFee() (fee float32, err error) {
	err = p.Model((*models.ServiceSettings)(nil)).
		Column("fee").
		Select(&fee)
	return
}
