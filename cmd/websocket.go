package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/plespurples/miniature-robot/internal/seats"
	"github.com/plespurples/miniature-robot/pkg/server"
)

func main() {
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
		server.Connections[counter] = c
		thisID := counter
		counter++

		var (
			mt  int
			msg []byte
			err error
		)

		// create a slice of strings from locked seats
		lockedStringList := []string{}
		for k := range seats.Locked {
			lockedStringList = append(lockedStringList, k)
		}

		// get current data
		cd := server.ResponseMessage{
			Event: "startstate",
			Data: seats.State{
				Reserved: []string{"A_2_2", "A_2_3", "A_2_1", "A_2_4", "V_60", "V_61"},
				Paid:     []string{"V_1", "V_2", "V_3", "V_4", "G_2_1", "G_2_2", "A_18_1", "A_18_2", "A_18_3", "A_18_4", "A_4_1", "A_4_2"},
				Locked:   lockedStringList,
			},
		}

		// run order timer
		go func() {
			time.Sleep(30 * time.Minute)

			// delete the locked seats and store their ids
			deletedSeats := []string{}
			for id, client := range seats.Locked {
				if client == thisID {
					// delete the seat
					deletedSeats = append(deletedSeats, id)
					delete(seats.Locked, id)

					// send unlocked message to all clients except this one
					server.BroadcastMessage(server.ResponseMessage{
						Event: "unlocked",
						Data:  id,
					}, thisID)
				}
			}

			// send the informative message to frontend
			server.SendMessage(c, server.ResponseMessage{
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
				log.Println("Error during reading the client message, aborting connection:", err.Error())
				break
			}

			// ok message received
			var sr seats.Request
			err := json.Unmarshal(msg, &sr)
			if err != nil {
				c.WriteMessage(mt, []byte("Hey, your JSON is invalid. Make it right!"))
				continue
			}

			// valid json, do the job
			switch sr.Action {
			case "lock":
				seats.HandleLock(c, sr, thisID)
			case "unlock":
				seats.HandleUnlock(c, sr, thisID)
			}
		}
	}))

	// start listening for requests
	log.Fatal(app.Listen(":3632"))
}
