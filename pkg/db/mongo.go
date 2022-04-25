package db

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/plespurples/miniature-robot/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client stores the mongo client for database requests
var Client *mongo.Client = nil

// MongoConnect creates the database connection. If it is not possible to
// establish one, an error is returned instead.
func MongoConnect(host string, db string, user string, pwd string) error {
	// create background context for the connection process
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// compose the mongo connection string
	options := options.Client().ApplyURI(fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/%s?retryWrites=true&w=majority", user, url.QueryEscape(pwd), host, db,
	))

	// try to connect and return the result immediately
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
