package postgres

import (
	"log"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/spf13/viper"

	"github.com/dmitriy-vas/node/models"
)

type postgres struct {
	*pg.DB
}

var tables = []interface{}{
	(*models.Bitcoin)(nil),
	(*models.BitcoinCash)(nil),
	(*models.Ethereum)(nil),
	(*models.Litecoin)(nil),
	(*models.Waves)(nil),
}
var Database *postgres

func init() {
	conn := pg.Connect(&pg.Options{
		Network:         viper.GetString("database.network"),
		Addr:            viper.GetString("database.address"),
		User:            viper.GetString("database.user"),
		Password:        viper.GetString("database.password"),
		Database:        viper.GetString("database.database"),
		ApplicationName: viper.GetString("database.application"),
	})
	options := &orm.CreateTableOptions{
		IfNotExists: true,
	}
	for _, table := range tables {
		if err := conn.CreateTable(table, options); err != nil {
			log.Panicf("Erorr with creating postgres tables: %v", err)
		}
	}
	Database = &postgres{conn}
}
