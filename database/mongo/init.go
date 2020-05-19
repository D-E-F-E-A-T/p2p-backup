package mongo

import (
	"context"
	"log"
	"time"

	"github.com/spf13/viper"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongo struct {
	*mongodb.Client
}

var Mongo *mongo

func init() {
	opts := options.Client()
	opts.SetAppName(viper.GetString("database.mongo.application"))
	opts.ApplyURI("mongodb://" + viper.GetString("database.mongo.address"))
	opts.SetAuth(options.Credential{
		Username: viper.GetString("database.mongo.user"),
		Password: viper.GetString("database.mongo.password"),
	})
	client, err := mongodb.NewClient(opts)
	if err != nil {
		log.Panicf("Error creating Mongo: %v", err)
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*20)
	if err := client.Connect(ctx); err != nil {
		log.Panicf("Error connecting to Mongo: %v", err)
	}
	Mongo = &mongo{client}
}
