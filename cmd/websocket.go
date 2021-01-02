package main

import (
	"encoding/json"
	"log"
	"time"

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
	Paid     []string `json:"paid"`
	Locked   []string `json:"locked"`
}

// ResponseMessage is a structure for sending server responses to the client
type ResponseMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

// connections is a map which contains all websocket connections
var connections = make(map[int]*websocket.Conn)

// counter is a counter for new connection ids
var counter = 0

// lockedSeats is a array of locked seat identifiers
var lockedSeats = make(map[string]time.Time)

// broadcastMessage sends a message to all clients, omitting
// the one which is specified in an omit param
func broadcastMessage(res ResponseMessage, omit int) {
	// send the change (event) to all connected clients
	for clientID, client := range connections {
		// omit current connection
		if clientID == omit {
			continue
		}

		// create the message string value from structure
		dStr, err := json.Marshal(res)
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
}

// handleLock locks the specified seat for lockerID and sends a message
// to all connected clients about the new locked seat (on success)
func handleLock(c *websocket.Conn, sr SeatRequest, lockerID int) {
	if _, ok := lockedSeats[sr.Seat]; ok {
		dStr, _ := json.Marshal(ResponseMessage{
			Event: "alreadyLocked",
			Data:  sr.Seat,
		})
		c.WriteMessage(1, dStr)
		return
	}

	// lock the seat for the specified amount of time
	lockedSeats[sr.Seat] = time.Now().Add(120 * time.Second)

	// create the message
	dStr, _ := json.Marshal(ResponseMessage{
		Event: "lockedForYou",
		Data:  sr.Seat,
	})

	// send messages back to the client
	c.WriteMessage(1, dStr)

	// send locked message to all clients
	broadcastMessage(ResponseMessage{
		Event: "locked",
		Data:  sr.Seat,
	}, lockerID)
}

// handleUnlock unlocks the specified seat and sends a message to all
// connected clients about the new unlocked seat (on success)
func handleUnlock(c *websocket.Conn, sr SeatRequest, lockerID int) {
	if _, ok := lockedSeats[sr.Seat]; !ok {
		dStr, _ := json.Marshal(ResponseMessage{
			Event: "alreadyUnlocked",
			Data:  sr.Seat,
		})
		c.WriteMessage(1, dStr)
		return
	}

	// unlock the seat
	delete(lockedSeats, sr.Seat)

	// create the message
	dStr, _ := json.Marshal(ResponseMessage{
		Event: "unlockedForYou",
		Data:  sr.Seat,
	})

	// send messages back to the client
	c.WriteMessage(1, dStr)

	// send locked message to all clients
	broadcastMessage(ResponseMessage{
		Event: "unlocked",
		Data:  sr.Seat,
	}, lockerID)
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
		// add this connection to the connections map
		connections[counter] = c
		thisID := counter
		counter++

		var (
			mt  int
			msg []byte
			err error
		)

		// create a slice of strings from locked seats
		lockedStringList := []string{}
		for k := range lockedSeats {
			lockedStringList = append(lockedStringList, k)
		}

		// get current data
		cd := ResponseMessage{
			Event: "startstate",
			Data: SeatsState{
				Reserved: []string{"A_2_2", "A_2_3", "A_2_1", "A_2_4", "V_60", "V_61"},
				Paid:     []string{"V_1", "V_2", "V_3", "V_4", "G_2_1", "G_2_2", "A_18_1", "A_18_2", "A_18_3", "A_18_4"},
				Locked:   lockedStringList,
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
			switch sr.Action {
			case "lock":
				handleLock(c, sr, thisID)
			case "unlock":
				handleUnlock(c, sr, thisID)
			}
		}
	}))

	// start listening for requests
	log.Fatal(app.Listen(":3632"))
}
