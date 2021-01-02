package seats

import (
	"encoding/json"

	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/server"
)

// HandleUnlock unlocks the specified seat and sends a message to all
// connected clients about the new unlocked seat (on success)
func HandleUnlock(c *websocket.Conn, sr Request, lockerID int) {
	if _, ok := Locked[sr.Seat]; !ok {
		dStr, _ := json.Marshal(server.ResponseMessage{
			Event: "alreadyUnlocked",
			Data:  sr.Seat,
		})
		c.WriteMessage(1, dStr)
		return
	}

	// unlock the seat
	delete(Locked, sr.Seat)

	// create the message
	dStr, _ := json.Marshal(server.ResponseMessage{
		Event: "unlockedForYou",
		Data:  sr.Seat,
	})

	// send messages back to the client
	c.WriteMessage(1, dStr)

	// send locked message to all clients
	server.BroadcastMessage(server.ResponseMessage{
		Event: "unlocked",
		Data:  sr.Seat,
	}, lockerID)
}
