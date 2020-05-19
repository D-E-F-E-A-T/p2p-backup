package postgres

import (
	"github.com/go-pg/pg/v9/orm"

	"github.com/dmitriy-vas/p2p/models"
)

func (p postgres) GetUserDealsAmount(id uint64) (count int, err error) {
	var m []*GetUserDealsModel
	return p.Model(&m).
		Relation("Offer").
		WhereOr("deal.user_id = ?", id).
		WhereOr("offer.user_id = ?", id).
		Count()
}

type GetUserDealsModel struct {
	User         *GetUserDealsUserModel `json:"user"`
	*models.Deal `pg:",inherit"`
	Offer        struct {
		tableName struct{}               `pg:"offers"`
		User      *GetUserDealsUserModel `json:"user"`
		*models.Offer
	} `json:"offer"`
}

type GetUserDealsUserModel struct {
	tableName struct{} `pg:"users"`
	ID        uint64   `json:"id" pg:",pk"`
	Login     string   `json:"login" pg:",unique,notnull"`
}

func (p postgres) GetUserDeals(id uint64, limit, offset int, sortMethod string) (count int, deals []*GetUserDealsModel, err error) {
	query := p.Model(&deals).
		Relation("User").
		Relation("Offer").
		Relation("Offer.User").
		WhereOr("deal.user_id = ?", id).
		WhereOr("offer.user_id = ?", id).
		Limit(limit).
		Offset(offset)
	switch sortMethod {
	case "New":
		query = query.Order("timestamp DESC")
	case "Old":
		query = query.Order("timestamp ASC")
	}
	count, err = query.SelectAndCount()
	return
}

type Deal struct {
	*models.Deal `pg:",inherit"`
	Argue        *models.Argue
	Chat         *models.Chat
	Offer        *models.Offer
}

func (p postgres) GetDeal(id uint64) (deal *Deal, err error) {
	deal = new(Deal)
	return deal, p.Model(deal).
		Relation("Offer").
		Relation("Argue").
		Relation("Chat", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("chat.first_user != 0").Where("chat.second_user != 0"), nil
		}).
		Where("deal.id = ?", id).
		Select()
}

func (p postgres) AddNewDeal(deal *models.Deal) error {
	return p.Insert(deal)
}

func (p postgres) SearchOfferDeals(id uint64) (count int, err error) {
	return p.Model((*models.Deal)(nil)).
		Where("offer_id = ?", id).
		Count()
}

func (p postgres) SetDealApproved(id uint64, place string, value bool) error {
	query := p.Model((*models.Deal)(nil)).
		Where("id = ?", id)
	switch place {
	case "From":
		query = query.Set("from_approved = ?", value)
	case "To":
		query = query.Set("to_approved = ?", value)
	}
	_, err := query.Update()
	return err
}

func (p postgres) SetDealAccepted(id uint64, value bool) error {
	_, err := p.Model((*models.Deal)(nil)).
		Where("id = ?", id).
		Set("accepted = ?", value).
		Update()
	return err
}

func (p postgres) DeleteDeal(id uint64) error {
	_, err := p.Model((*models.Deal)(nil)).
		Where("id = ?", id).
		Delete()
	return err
}

func (p postgres) SetDealFinished(id uint64, value bool) error {
	_, err := p.Model((*models.Deal)(nil)).
		Where("id = ?", id).
		Where("finished = ?", !value).
		Set("finished = ?", value).
		Update()
	return err
}

func (p postgres) SetFixedAndFee(id uint64, amount float32, fee float32) error {
	_, err := p.Model((*models.Deal)(nil)).
		Where("id = ?", id).
		Set("fixed_amount = ?", amount).
		Update()
	return err
}

func (p postgres) SetDealWithMessages(id uint64) error {
	_, err := p.Model((*models.Deal)(nil)).
		Where("id = ?", id).
		Set("with_messages = ?", true).
		Update()
	return err
}

func (p postgres) GetReserveDeals(limit int, offset int) (count int, deals []*models.Deal, err error) {
	query := p.Model(&deals).
		Where("fixed_amount IS NOT NULL").
		Limit(limit).
		Offset(offset).
		Order("timestamp DESC")
	count, err = query.SelectAndCount()
	return
}
