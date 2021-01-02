package seats

import (
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/server"
)

// HandleLock locks the specified seat for lockerID and sends a message
// to all connected clients about the new locked seat (on success)
func HandleLock(c *websocket.Conn, sr Request, lockerID int) {
	if _, ok := Locked[sr.Seat]; ok {
		server.SendMessage(c, server.ResponseMessage{
			Event: "alreadyLocked",
			Data:  sr.Seat,
		})
		return
	}

	// lock the seat for the specified amount of time
	Locked[sr.Seat] = time.Now().Add(120 * time.Second)

	// send success message to the locking client
	server.SendMessage(c, server.ResponseMessage{
		Event: "lockedForYou",
		Data:  sr.Seat,
	})

	// send locked message to all clients
	server.BroadcastMessage(server.ResponseMessage{
		Event: "locked",
		Data:  sr.Seat,
	}, lockerID)
}
