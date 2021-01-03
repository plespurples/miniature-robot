package seats

import (
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

// HandleUnlock unlocks the specified seat and sends a message to all
// connected clients about the new unlocked seat (on success)
func HandleUnlock(c *websocket.Conn, sr Request, unlockerID int) {
	// if the seat is locked for other person than unlockerID, send
	// corresponding message to the client
	if Locked[sr.Seat] != unlockerID {
		wssrv.SendMessage(c, wssrv.ResponseMessage{
			Event: "notYours",
			Data:  sr.Seat,
		})
		return
	}

	// if the seat is already unlocked, stop the job
	if _, ok := Locked[sr.Seat]; !ok {
		wssrv.SendMessage(c, wssrv.ResponseMessage{
			Event: "alreadyUnlocked",
			Data:  sr.Seat,
		})
		return
	}

	// unlock the seat
	delete(Locked, sr.Seat)

	// send the message
	wssrv.SendMessage(c, wssrv.ResponseMessage{
		Event: "unlockedForYou",
		Data:  sr.Seat,
	})

	// send locked message to all clients
	wssrv.BroadcastMessage(wssrv.ResponseMessage{
		Event: "unlocked",
		Data:  sr.Seat,
	}, unlockerID)
}
