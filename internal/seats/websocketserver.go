package seats

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/pkg/config"
	"github.com/plespurples/miniature-robot/pkg/wssrv"
)

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

		var (
			mt  int
			msg []byte
			err error
		)

		// create a slice of strings from locked seats
		lockedStringList := []string{}
		for k := range State.Locked {
			lockedStringList = append(lockedStringList, k)
		}

		// create a slice of strings from reserved seats
		reservedStringList := []string{}
		for k := range State.Reserved {
			reservedStringList = append(reservedStringList, k)
		}

		// create a slice of strings from paid seats
		paidStringList := []string{}
		for k := range State.Paid {
			paidStringList = append(paidStringList, k)
		}

		// get current data
		cd := wssrv.ResponseMessage{
			Event: "startstate",
			Data: map[string][]string{
				"reserved": reservedStringList,
				"paid":     paidStringList,
				"locked":   lockedStringList,
			},
		}

		// run order timer
		go func() {
			time.Sleep(30 * time.Minute)

			// delete the locked seats and store their ids
			deletedSeats := []string{}
			for id, client := range State.Locked {
				if client == thisID {
					// delete the seat
					deletedSeats = append(deletedSeats, id)
					delete(State.Locked, id)

					// send unlocked message to all clients except this one
					wssrv.BroadcastMessage(wssrv.ResponseMessage{
						Event: "unlocked",
						Data:  id,
					}, thisID)
				}
			}

			// send the informative message to frontend
			wssrv.SendMessage(c, wssrv.ResponseMessage{
				Event: "deleted",
				Data:  deletedSeats,
			})
		}()

		// marshal the data to json string
		currentStateString, err := json.Marshal(cd)
		if err != nil {
			log.Println("Cannot marshal the current data", cd)
		}

		// this happen when client is connected to the server
		if err = c.WriteMessage(1, currentStateString); err != nil {
			log.Println("Error while sending message:", err.Error())
		}

		// this will happen on every message/connection
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				// unlock all seats and close the connection
				UnlockAll(thisID)
				c.Close()
				delete(wssrv.Connections, thisID)

				log.Println("Aborted connection:", err.Error())
				break
			}

			// ok message received
			var sr Request
			err := json.Unmarshal(msg, &sr)
			if err != nil {
				c.WriteMessage(mt, []byte("Hey, your JSON is invalid. Make it right!"))
				continue
			}

			// valid json, do the job
			switch sr.Action {
			case "lock":
				HandleLock(c, sr, thisID)
			case "unlock":
				HandleUnlock(c, sr, thisID)
			case "reserve":
				HandleReserve(c, sr, thisID)
			}
		}
	}))

	// start listening for requests
	log.Fatal(app.Listen(":" + config.Data.Net.Port))
}
