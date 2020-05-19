package postgres

import (
	"context"
	"log"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/spf13/viper"

	"github.com/dmitriy-vas/p2p/models"
)

type postgres struct {
	*pg.DB
}

type postgresDebugger struct{}

func (d postgresDebugger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d postgresDebugger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	log.Println(q.FormattedQuery())
	return nil
}

var (
	tables = []interface{}{
		(*models.Argue)(nil),
		(*models.Attachment)(nil),
		(*models.Auth)(nil),
		(*models.Ban)(nil),
		(*models.Chat)(nil),
		(*models.Comment)(nil),
		(*models.Deal)(nil),
		(*models.Entry)(nil),
		(*models.Message)(nil),
		(*models.Notification)(nil),
		(*models.Offer)(nil),
		(*models.Relationship)(nil),
		(*models.Restore)(nil),
		(*models.ServiceSettings)(nil),
		(*models.Session)(nil),
		(*models.Settings)(nil),
		(*models.Statistic)(nil),
		(*models.Transaction)(nil),
		(*models.User)(nil),
		(*models.Wallet)(nil),
	}
	Postgres *postgres
)

const (
	LanguageEnglish = "en"
	LanguageRussian = "ru"
)

func init() {
	conn := pg.Connect(&pg.Options{
		Network:         viper.GetString("database.postgres.network"),
		Addr:            viper.GetString("database.postgres.address"),
		User:            viper.GetString("database.postgres.user"),
		Password:        viper.GetString("database.postgres.password"),
		Database:        viper.GetString("database.postgres.database"),
		ApplicationName: viper.GetString("database.postgres.application"),
	})
	options := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	for _, table := range tables {
		if err := conn.CreateTable(table, options); err != nil {
			log.Panicf("Erorr with creating postgres tables: %v", err)
		}
	}
	conn.AddQueryHook(postgresDebugger{})
	Postgres = &postgres{conn}
	go Postgres.CleanExpired()
}

func (p postgres) CleanExpired() {
	for {
		if _, err := p.Model((*models.Restore)(nil)).
			Where("expires < ?", time.Now()).
			ForceDelete(); err != nil {
			log.Printf("Error with deleting expired restore: %v", err)
		}
		if _, err := p.Model((*models.Session)(nil)).
			Where("expires < ?", time.Now()).
			ForceDelete(); err != nil {
			log.Printf("Error with deleting expired session: %v", err)
		}
		// TODO delete expired deals
		//var m []*GetUserDealsModel
		//Postgres.Model(&m).
		//	Relation("Offer").
		//	Where("deal.accepted = ?", false).
		//	Where("deal.timestamp < current_timestamp + (offer.cancel_time || ' minutes')::interval", ).
		//	Select()
		//if len(m) > 0 {
		//	var deals []*models.Deal
		//	for _, deal := range m {
		//		deals = append(deals, deal.Deal)
		//	}
		//	if _, err := Postgres.Model(&deals).WherePK().Delete(); err != nil {
		//		log.Printf("Error with deleting expired deals: %v", err)
		//	}
		//}
		if _, err := p.Model((*models.Auth)(nil)).
			Where("expires < ?", time.Now()).
			ForceDelete(); err != nil {
			log.Printf("Error with deleting expired auth codes: %v", err)
		}
		if _, err := p.Model((*models.Ban)(nil)).
			Where("expires < ?", time.Now()).
			ForceDelete(); err != nil {
			log.Printf("Error with deleting expired bans: %v", err)
		}
		time.Sleep(time.Minute * 30)
	}
}
