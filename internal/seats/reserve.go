package seats

import (
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

// HandleReserve ...
func HandleReserve(c *websocket.Conn, sr Request, lockerID int) {
	if sr.AuthorizationString != "TestString123" {
		wssrv.SendMessage(c, wssrv.ResponseMessage{
			Event: "unauthorized",
			Data:  "You are unauthorized to do this.",
		})
		return
	}

	// update the list of reserved seats
	end := time.Now().Add(365 * 24 * time.Hour)
	Reserved[sr.Seat] = &end

	// send locked message to all clients
	wssrv.BroadcastMessage(wssrv.ResponseMessage{
		Event: "reserved",
		Data:  sr.Seat,
	}, lockerID)
}
