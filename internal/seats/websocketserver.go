package seats

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/config"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

// getStartingData gets current seat statuses on new websocket
// connection initialization. It is returned as map of strings
// with a slice-of-strings value type which can be immediately
// passed to the ResponseMessage Data field.
func getStartingData() map[string][]string {
	lockedStringList := []string{}
	for k := range State.Locked {
		lockedStringList = append(lockedStringList, k)
	}

	reservedStringList := []string{}
	for k := range State.Reserved {
		reservedStringList = append(reservedStringList, k)
	}

	paidStringList := []string{}
	for k := range State.Paid {
		paidStringList = append(paidStringList, k)
	}

	return map[string][]string{
		"reserved": reservedStringList,
		"paid":     paidStringList,
		"locked":   lockedStringList,
	}
}

// handleMessage handles one message that is received by our websocket
// server. It unmarshal the json message into the request struct and call
// appropriate function depending to the specified action from the message.
func handleMessage(msg []byte, c *websocket.Conn, id int) {
	// ok message received
	var sr Request
	err := json.Unmarshal(msg, &sr)
	if err != nil {
		c.WriteMessage(1, []byte("Hey, your JSON is invalid. Make it right!"))
		return
	}

	// valid json, do the job
	switch sr.Action {
	case "lock":
		HandleLock(c, sr, id)
	case "unlock":
		HandleUnlock(c, sr, id)
	case "reserve":
		HandleReserve(c, sr, id)
	case "unreserve":
		HandleUnreserve(c, sr, id)
	case "revoke":
		HandleRevoke(c, sr, id)
	case "pay":
		HandlePay(c, sr, id)
	}
}

// RunWebsocketServer starts the websocket server which is used to handle
// all seat clicks. It locks or unlocks the seat on click on the website
// and prevents locking too many seats for one user (session), also it
// prevents having locked seats for a lot of time by setting time limit
func RunWebsocketServer() {
	// connection counter
	counter := 0

	// create new server
	app := fiber.New()

	// ensure upgrading to the websocket protocol
	app.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// websocket endpoint
	app.Get("/", websocket.New(func(c *websocket.Conn) {
		// add this connection to the connections map
		wssrv.Connections[counter] = c
		thisID := counter
		counter++

		// get current data
		cd := wssrv.ResponseMessage{
			Event: "startstate",
			Data:  getStartingData(),
		}

		// marshal the data to json string and send them to client
		currentStateString, _ := json.Marshal(cd)
		if err := c.WriteMessage(1, currentStateString); err != nil {
			log.Println("Error while sending message:", err.Error())
		}

		// run order creation timer, when this timer expires, all the
		// locked places of this connection will be unlocked to other people
		// todo: implement something like this for web users
		// go func(c *websocket.Conn, id int) {
		// 	time.Sleep(1 * time.Minute)

		// 	// send the informative message to frontend and unlock all seats
		// 	wssrv.SendMessage(c, wssrv.ResponseMessage{
		// 		Event: "deleted",
		// 		Data:  GetLocked(id),
		// 	})
		// 	UnlockAll(id)
		// }(c, thisID)

		// this will happen on every message/connection
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				// unlock all seats and close the connection
				UnlockAll(thisID)
				c.Close()
				delete(wssrv.Connections, thisID)

				log.Println("Aborted connection:", err.Error())
				break
			}

			// handle the message
			handleMessage(msg, c, thisID)
		}
	}))

	// start listening for requests
	log.Fatal(app.Listen(":" + config.Data.Net.Port))
}
