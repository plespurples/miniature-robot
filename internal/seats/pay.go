package seats

import (
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/config"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

// HandlePay handles the requests to set a seat to be busy. It
// required an authorization by providing the authorization string.
func HandlePay(c *websocket.Conn, sr Request, lockerID int) {
	if sr.AuthorizationString != config.Data.Security.AuthorizationString {
		wssrv.SendMessage(c, wssrv.ResponseMessage{
			Event: "unauthorized",
			Data:  "You are unauthorized to do this.",
		})
		return
	}

	// update the list of reserved seats
	delete(State.Reserved, sr.Seat)
	State.Paid[sr.Seat] = struct{}{}

	// send locked message to all clients
	wssrv.BroadcastMessage(wssrv.ResponseMessage{
		Event: "paid",
		Data:  sr.Seat,
	}, lockerID)
}
