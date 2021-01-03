package db

import (
	"context"
	"net/url"
	"time"

	"github.com/plespurples/miniature-robot/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client stores the mongo client for database requests
var Client *mongo.Client = nil

// MongoConnect creates the database connection
func MongoConnect() error {
	// try to connect to the database
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	options := options.Client().ApplyURI(
		"mongodb://" +
			config.Data.Database.User +
			":" +
			url.QueryEscape(config.Data.Database.Password) +
			"@" +
			config.Data.Database.Host +
			":" +
			config.Data.Database.Port +
			"/?authSource=" +
			config.Data.Database.Name +
			"&connect=direct")

	c, err := mongo.Connect(ctx, options)
	if err != nil {
		return err
	}
	Client = c
	return nil
}

// Collection returns the specified collection
func Collection(c string) *mongo.Collection {
	return Client.Database(config.Data.Database.Name).Collection(c)
}
