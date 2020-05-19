package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) AddNewArgue(argue *models.Argue) error {
	return p.Insert(argue)
}

func (p postgres) SearchArgue(id uint64) (argue *models.Argue, err error) {
	argue = new(models.Argue)
	return argue, p.Model(argue).
		Where("id = ?", id).
		Select()
}

func (p postgres) DeleteArgue(id uint64) error {
	_, err := p.Model((*models.Argue)(nil)).
		Where("id = ?", id).
		Delete()
	return err
}

func (p postgres) FinishArgue(id uint64) error {
	_, err := p.Model((*models.Argue)(nil)).
		Where("id = ?", id).
		Set("finished = ?", true).
		Update()
	return err
}
