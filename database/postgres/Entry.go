package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) AddNewEntry(entry *models.Entry) error {
	return p.Insert(entry)
}

func (p postgres) SearchUserEntries(id uint64, limit, offset int) (count int, entries []*models.Entry, err error) {
	count, err = p.Model(&entries).
		Where("id = ?", id).
		Limit(limit).
		Offset(offset).
		SelectAndCount()
	return
}
