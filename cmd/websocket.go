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

// SeatsState is a total state for all seats in the whole hall (in all rooms)
type SeatsState struct {
	Reserved []string `json:"reserved"`
	Busy     []string `json:"busy"`
	Locked   []string `json:"locked"`
}

// ResponseMessage is a structure for sending server responses to the client
type ResponseMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// connections is a map which contains all websocket connections
var connections map[int]*websocket.Conn = make(map[int]*websocket.Conn)

// counter is a counter for new connection ids
var counter int = 0

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
		// add this connection to the connections map
		connections[counter] = c
		counter++

		var (
			mt  int
			msg []byte
			err error
		)

		// get current data
		cd := ResponseMessage{
			Event: "startstate",
			Data: SeatsState{
				Reserved: []string{"A_2_2", "A_2_3", "A_2_1", "A_2_4", "V_60", "V_61"},
				Busy:     []string{"V_1", "V_2", "V_3", "V_4", "G_2_1", "G_2_2", "A_18_1", "A_18_2", "A_18_3", "A_18_4"},
				Locked:   []string{"A_10_1", "A_10_2"},
			},
		}

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
			var sr SeatRequest
			err := json.Unmarshal(msg, &sr)
			if err != nil {
				c.WriteMessage(mt, []byte("Hey, your JSON is invalid. Make it right!"))
				continue
			}

			// valid json, do the job
			fmt.Println(sr.Action)
			fmt.Println(sr.Seat)

			// send the change (event) to all connected clients
			for _, client := range connections {
				// create the message string value from structure
				dStr, err := json.Marshal(ResponseMessage{
					Event: sr.Action + "ed",
					Data:  sr.Seat,
				})
				if err != nil {
					log.Println("Error:", err.Error())
					continue
				}

				// send the message
				err = client.WriteMessage(1, dStr)
				if err != nil {
					log.Println("Error:", err.Error())
					continue
				}
			}

			// send messages back to the client
			smsg := []byte("Super")
			if err = c.WriteMessage(mt, smsg); err != nil {
				log.Println("Error while sending message:", err.Error())
				break
			}
		}
	}))

	// start listening for requests
	log.Fatal(app.Listen(":3632"))
}
