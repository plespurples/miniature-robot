package seats

import (
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/config"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

// HandleRevoke handles the requests to revoke a seat. It
// required an authorization by providing the authorization string.
func HandleRevoke(c *websocket.Conn, sr Request, lockerID int) {
	if sr.AuthorizationString != config.Data.Security.AuthorizationString {
		wssrv.SendMessage(c, wssrv.ResponseMessage{
			Event: "unauthorized",
			Data:  "You are unauthorized to do this.",
		})
		return
	}

	// update the list of reserved seats
	delete(State.Reserved, sr.Seat)
	delete(State.Paid, sr.Seat)

	// send a message to all clients
	wssrv.BroadcastMessage(wssrv.ResponseMessage{
		Event: "revoked",
		Data:  sr.Seat,
	}, lockerID)
}
