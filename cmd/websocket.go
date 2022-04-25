package main

import (
	"os"

	"github.com/plespurples/miniature-robot/internal/seats"
	"github.com/plespurples/miniature-robot/pkg/config"
	"github.com/plespurples/miniature-robot/pkg/db"
)

func main() {
	config.Load(os.Getenv("CONFIG"))
	err := db.MongoConnect(
		config.Data.Database.Host,
		config.Data.Database.Name,
		config.Data.Database.User,
		config.Data.Database.Password,
	)
	if err != nil {
		panic(err)
	}
	seats.SetCurrentStatus()
	seats.RunWebsocketServer()
}
