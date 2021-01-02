package server

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
)

// ResponseMessage is a structure for sending server responses to the client
type ResponseMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// Connections is a map which contains all websocket connections
var Connections = make(map[int]*websocket.Conn)

// BroadcastMessage sends a message to all clients, omitting
// the one which is specified in an omit param
func BroadcastMessage(res ResponseMessage, omit int) {
	// send the change (event) to all connected clients
	for clientID, client := range Connections {
		// omit current connection
		if clientID == omit {
			continue
		}

		// create the message string value from structure
		dStr, err := json.Marshal(res)
		if err != nil {
			log.Println("Error:", err.Error())
			continue
		}

		// send the message
		err = client.WriteMessage(1, dStr)
		if err != nil {
			log.Println("Error:", err.Error())
			continue
		}
	}
}
