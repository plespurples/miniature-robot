package wssrv

import (
	"encoding/json"

	"github.com/gofiber/websocket/v2"
)

// ResponseMessage is a structure for sending server responses to the client
type ResponseMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// Connections is a map which contains all websocket connections
var Connections = make(map[int]*websocket.Conn)

// SendMessage sends message to the specified connection
func SendMessage(c *websocket.Conn, rm ResponseMessage) {
	msg, _ := json.Marshal(rm)
	c.WriteMessage(1, msg)
}

// BroadcastMessage sends a message to all clients, omitting
// the one which is specified in an omit param
func BroadcastMessage(res ResponseMessage, omit int) {
	// send the change (event) to all connected clients
	for clientID, client := range Connections {
		// omit the specified client
		if clientID == omit {
			continue
		}

		// let's send the message
		SendMessage(client, res)
	}
}
