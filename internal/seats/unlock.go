package seats

import (
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

// Unlock unlocks one particular seat which is defined by the
// seat parameter. It also broadcasts a message about the unlock
// process to all listening connections (excepting the ommited one).
func Unlock(seat string, omit int) {
	delete(Locked, seat)
	wssrv.BroadcastMessage(wssrv.ResponseMessage{
		Event: "unlocked",
		Data:  seat,
	}, omit)
}

// UnlockAll unlocks all seats that belongs to the locker ID. From
// the Unlock function which is called for unlocking every seat, it
// also broadcasts the message about the unlocking to all connections.
func UnlockAll(lockerID int) {
	for seat, lid := range Locked {
		if lid == lockerID {
			Unlock(seat, lockerID)
		}
	}
}

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
	Unlock(sr.Seat, unlockerID)

	// send the success message to the unlocking connection
	wssrv.SendMessage(c, wssrv.ResponseMessage{
		Event: "unlockedForYou",
		Data:  sr.Seat,
	})
}
