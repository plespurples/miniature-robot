package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SeatRequest is a structure which represents a payload from websocket
// informing backend about one specific action (clicked or unclicked seat)
type SeatRequest struct {
	Action string `json:"action"`
	Seat   string `json:"seat"`
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
		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)

		// this happen when client is connected to the server
		fmt.Println("New connection established")

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
	log.Fatal(app.Listen(":3000"))
}
