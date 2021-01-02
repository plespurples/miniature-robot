package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SeatRequest is a structure which represents a payload from websocket
// informing backend about one specific action (locked or unlocked seat)
type SeatRequest struct {
	Action string `json:"action"`
	Seat   string `json:"seat"`
}

type SeatsState struct {
	Reserved []string `json:"reserved"`
	Busy     []string `json:"busy"`
	Locked   []string `json:"locked"`
}

func main() {
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
		var (
			mt  int
			msg []byte
			err error
		)

		// get current data
		cd := SeatsState{
			Reserved: []string{"A_2_2", "A_2_3", "A_2_1", "A_2_4", "V_60", "V_61"},
			Busy:     []string{"V_1", "V_2", "V_3", "V_4", "G_2_1", "G_2_2", "A_18_1", "A_18_2", "A_18_3", "A_18_4"},
			Locked:   []string{"A_10_1", "A_10_2"},
		}

		// marshal the data to json string
		currentStateString, err := json.Marshal(cd)
		if err != nil {
			log.Println("Cannot marshal the current data", cd)
		}

		// this happen when client is connected to the server
		if err = c.WriteMessage(1, currentStateString); err != nil {
			log.Println("Error while sending message:", err)
		}

		// this will happen on every message/connection
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Println("Error during reading the client message, aborting connection:", err)
				break
			}

			// ok message received
			var sr SeatRequest
			err := json.Unmarshal(msg, &sr)
			if err != nil {
				c.WriteMessage(mt, []byte("Hey, your JSON is invalid. Make it right!"))
				continue
			}

			// valid json, do the job
			fmt.Println(sr.Action)
			fmt.Println(sr.Seat)

			// send messages back to the client
			smsg := []byte("Super")
			if err = c.WriteMessage(mt, smsg); err != nil {
				log.Println("Error while sending message:", err)
				break
			}
		}
	}))

	// start listening for requests
	log.Fatal(app.Listen(":3632"))
}
