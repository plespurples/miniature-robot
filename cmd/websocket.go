package main

import (
	"github.com/plespurples/miniature-robot/internal/seats"
	"github.com/plespurples/miniature-robot/pkg/config"
)

func main() {
	config.Load()
	seats.RunWebsocketServer()
}
