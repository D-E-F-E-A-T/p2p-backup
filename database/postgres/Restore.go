package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) AddNewRestore(restore *models.Restore) error {
	return p.Insert(restore)
}

func (p postgres) SearchRestoreToken(token string) (restore *models.Restore, err error) {
	restore = new(models.Restore)
	return restore, p.Model(restore).
		Where("token = ?", token).
		Select()
}

func (p postgres) DeleteRestoreToken(token string) error {
	_, err := p.Model((*models.Restore)(nil)).
		Where("token = ?", token).
		ForceDelete()
	return err
}
