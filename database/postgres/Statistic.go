package postgres

import (
	"time"

	"github.com/go-pg/pg/v9"

	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) SearchUserStatistic(id uint64) (statistic *models.Statistic, err error) {
	statistic = new(models.Statistic)
	err = p.Model(statistic).
		Where("id = ?", id).
		Select()
	if err == pg.ErrNoRows {
		statistic.ID = id
		statistic.UserID = id
		statistic.Created = time.Now()
		err = p.Insert(statistic)
	}
	return
}

func (p postgres) UpdateLastSeen(id uint64, t time.Time) error {
	_, err := p.Model((*models.Statistic)(nil)).
		Where("id = ?", id).
		Set("last_seen = ?", t).
		Update()
	return err
}

func (p postgres) UpdateEmailVerified(id uint64, value bool) error {
	statistic, err := p.SearchUserStatistic(id)
	if err != nil {
		return err
	}
	_, err = p.Model(statistic).
		Where("id = ?", id).
		Set("email_verified = ?", value).
		Update()
	return err
}

func (p postgres) IncrementFinishedDeals(id uint64) error {
	_, err := p.Model((*models.Statistic)(nil)).
		Where("id = ?", id).
		Set("finished_deals = COALESCE(finished_deals,0) + 1").
		Update()
	return err
}

func (p postgres) IncreaseUserRating(id uint64, amount uint8) error {
	_, err := p.Model((*models.Statistic)(nil)).
		Where("id = ?", id).
		Set("rating = COALESCE(rating,0) + ?", amount).
		Update()
	return err
}
