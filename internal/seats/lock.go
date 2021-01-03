package seats

import (
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

// HandleLock locks the specified seat for lockerID and sends a message
// to all connected clients about the new locked seat (on success).
func HandleLock(c *websocket.Conn, sr Request, lockerID int) {
	if _, ok := State.Locked[sr.Seat]; ok {
		wssrv.SendMessage(c, wssrv.ResponseMessage{
			Event: "alreadyLocked",
			Data:  sr.Seat,
		})
		return
	}

	// lock the seat for the specified amount of time
	State.Locked[sr.Seat] = lockerID

	// send success message to the locking client
	wssrv.SendMessage(c, wssrv.ResponseMessage{
		Event: "lockedForYou",
		Data:  sr.Seat,
	})

	// send locked message to all clients
	wssrv.BroadcastMessage(wssrv.ResponseMessage{
		Event: "locked",
		Data:  sr.Seat,
	}, lockerID)
}
