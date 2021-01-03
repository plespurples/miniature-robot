package seats

import (
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/config"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

// HandleReserve handles the requests to set a seat to be reserved. It
// required an authorization by providing the authorization string.
func HandleReserve(c *websocket.Conn, sr Request, lockerID int) {
	if sr.AuthorizationString != config.Data.Security.AuthorizationString {
		wssrv.SendMessage(c, wssrv.ResponseMessage{
			Event: "unauthorized",
			Data:  "You are unauthorized to do this.",
		})
		return
	}

	// update the list of reserved seats
	end := time.Now().Add(365 * 24 * time.Hour) // todo handle the end by another way
	State.Reserved[sr.Seat] = &end

	// send locked message to all clients
	wssrv.BroadcastMessage(wssrv.ResponseMessage{
		Event: "reserved",
		Data:  sr.Seat,
	}, lockerID)
}
