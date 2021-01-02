package seats

import (
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/server"
)

// HandleUnlock unlocks the specified seat and sends a message to all
// connected clients about the new unlocked seat (on success)
func HandleUnlock(c *websocket.Conn, sr Request, lockerID int) {
	if _, ok := Locked[sr.Seat]; !ok {
		server.SendMessage(c, server.ResponseMessage{
			Event: "alreadyUnlocked",
			Data:  sr.Seat,
		})
		return
	}

	// unlock the seat
	delete(Locked, sr.Seat)

	// send the message
	server.SendMessage(c, server.ResponseMessage{
		Event: "unlockedForYou",
		Data:  sr.Seat,
	})

	// send locked message to all clients
	server.BroadcastMessage(server.ResponseMessage{
		Event: "unlocked",
		Data:  sr.Seat,
	}, lockerID)
}
