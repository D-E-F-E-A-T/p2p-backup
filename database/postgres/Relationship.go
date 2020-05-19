package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) GetUserRelationship(id uint64, target uint64) (relationship *models.Relationship, err error) {
	relationship = new(models.Relationship)
	return relationship, p.Model(relationship).
		Where("user_id = ?", id).
		Where("id = ?", target).
		Select()
}
