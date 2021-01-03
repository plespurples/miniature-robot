package main

import (
	"github.com/plespurples/miniature-robot/internal/seats"
	"github.com/plespurples/miniature-robot/pkg/config"
	"github.com/plespurples/miniature-robot/pkg/db"
)

func main() {
	config.Load()
	db.MongoConnect()
	seats.SetCurrentStatus()
	seats.RunWebsocketServer()
}
