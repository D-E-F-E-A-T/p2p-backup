package postgres

import (
	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) SearchUserNotifications(id uint64, limit, offset int) (count int, notifications []*models.Notification, err error) {
	count, err = p.Model(&notifications).
		Where("id = ?", id).
		Limit(limit).
		Offset(offset).
		SelectAndCount()
	return
}

func (p postgres) AddNewNotification(notification *models.Notification) error {
	return p.Insert(notification)
}
