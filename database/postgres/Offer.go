package postgres

import (
	"time"

	"github.com/go-pg/pg/v9/orm"

	"github.com/dmitriy-vas/p2p/models"
)

type SearchOffersUser struct {
	tableName struct{} `pg:"users"`
	ID        uint64   `json:"id" pg:",pk"`
	Login     string   `json:"login" pg:",unique,notnull"`
}

type SearchOffersStatistic struct {
	tableName     struct{}  `pg:"statistics"`
	ID            uint64    `json:"id" pg:",pk"`
	PhoneVerified bool      `json:"phone_verified" pg:",use_zero"`
	DocsVerified  bool      `json:"docs_verified" pg:",use_zero"`
	LastSeen      time.Time `json:"last_seen"`
	Created       time.Time `json:"created"`
	FinishedDeals int       `json:"finished_deals"`
	DealsVolume   uint      `json:"deals_volume"`
	Rating        int       `json:"rating"`
}

type SearchOffersSettings struct {
	tableName struct{} `pg:"settings"`
	ID        uint64   `json:"id" pg:",pk"`
	Timezone  string   `json:"timezone"`
}

type Offer struct {
	*models.Offer `json:"offer" pg:",inherit"`
	User          *SearchOffersUser      `json:"user"`
	Statistic     *SearchOffersStatistic `json:"statistic"`
	Settings      *SearchOffersSettings  `json:"settings"`
}

func (p postgres) SearchOffers(offerMethod string,
	cryptoCurrency uint8,
	currency uint8,
	location string,
	search string,
	amount float32,
	sortMethod string,
	warranty bool,
	active bool,
	provider uint8,
	user uint64,
	limit int,
	offset int) (count int, offers []*Offer, err error) {
	query := p.Model(&offers).
		Relation("Statistic").
		Relation("Settings").
		Relation("User", func(q *orm.Query) (*orm.Query, error) {
			return q, nil
		}).
		Limit(limit).
		Offset(offset)
	if active {
		query = query.Where("active = ?", true)
	}
	if warranty {
		query = query.Where("warranty = ?", warranty)
	}
	if user != 0 {
		query = query.Where("offer.user_id = ?", user)
	}
	if provider != 0 {
		query = query.Where("provider = ?", provider)
	}
	switch sortMethod {
	case "Popularity":
		query = query.Order("rating DESC")
	case "New":
		query = query.Order("timestamp DESC")
	case "Old":
		query = query.Order("timestamp ASC")
	case "Cost":
		query = query.Order("cost ASC")
	}
	if location != "" {
		query = query.Where("location = ?", location)
	}
	if search != "" {
		query = query.Where("offer.title ~ ?", search)
	}
	if amount != 0 {
		query = query.WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
			return q.Where("minimal < ?", amount).
				Where("maximal > ?", amount), nil
		}).WhereOrGroup(func(q *orm.Query) (*orm.Query, error) {
			return q.Where("minimal = 0").
				Where("maximal = 0"), nil
		})
	}
	if cryptoCurrency != 0 {
		switch offerMethod {
		case "SellOrBuy", "":
			query = query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				return q.WhereOr("want_currency = ?", cryptoCurrency).
					WhereOr("have_currency = ?", cryptoCurrency), nil
			})
		case "Sell":
			query = query.Where("want_currency = ?", cryptoCurrency)
		case "Buy":
			query = query.Where("have_currency = ?", cryptoCurrency)
		}
	}
	if currency != 0 {
		switch offerMethod {
		case "SellOrBuy", "":
			query = query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
				return q.WhereOr("want_currency = ?", currency).
					WhereOr("have_currency = ?", currency), nil
			})
		case "Sell":
			query = query.Where("have_currency = ?", currency)
		case "Buy":
			query = query.Where("want_currency = ?", currency)
		}
	}
	count, err = query.SelectAndCount()
	return
}

func (p postgres) SearchOffersOld(have, want, provider string, limit, offset int) (count int, offers []*Offer, err error) {
	query := p.Model(&offers).
		Relation("User").
		Relation("Statistic").
		Relation("Settings").
		Limit(limit).
		Offset(offset)
	if provider != "" {
		query = query.Relation("Providers", func(q *orm.Query) (*orm.Query, error) {
			return q.Where("name = ?", provider), nil
		})
	}
	if have != "" {
		query = query.Where("offer.have_currency = ?", have)
	}
	if want != "" {
		query = query.Where("offer.want_currency = ?", want)
	}
	count, err = query.SelectAndCount()
	return
}

func (p postgres) SearchOffer(id uint64) (offer *Offer, err error) {
	offer = new(Offer)
	return offer, p.Model(offer).
		Relation("Statistic").
		Relation("Settings").
		Relation("User", func(q *orm.Query) (*orm.Query, error) {
			return q, nil
		}).
		Where("offer.id = ?", id).
		Select()
}

func (p postgres) GetUserOffersAmount(id uint64) (count int, err error) {
	return p.Model((*models.Offer)(nil)).
		Where("user_id = ?", id).
		Count()
}

func (p postgres) AddNewOffer(offer *models.Offer) error {
	return p.Insert(offer)
}

func (p postgres) GetOffer(id uint64) (offer *models.Offer, err error) {
	offer = new(models.Offer)
	return offer, p.Model(offer).
		Where("id = ?", id).
		Select()
}

func (p postgres) UpdateOffer(offer *models.Offer) error {
	return p.Update(offer)
}

func (p postgres) DeleteOffer(id uint64) error {
	_, err := p.Model((*models.Offer)(nil)).
		Where("id = ?", id).
		Delete()
	return err
}

func (p postgres) ToggleOffer(id uint64, user uint64) error {
	_, err := p.Model((*models.Offer)(nil)).
		Where("id = ?", id).
		Where("user_id = ?", user).
		Set("active = NOT active").
		Update()
	return err
}
