package seats

import (
	"context"
	"log"
	"time"

	"github.com/plespurples/miniature-robot/pkg/config"
	"github.com/plespurples/miniature-robot/pkg/db"
	"go.mongodb.org/mongo-driver/bson"
)

// Request is a structure which represents a payload from websocket
// informing backend about one specific action (locked or unlocked seat).
type Request struct {
	Action              string `json:"action"`
	Seat                string `json:"seat"`
	AuthorizationString string `json:"authorizationString"`
}

// state represents the data about all seats
type state struct {
	// Locked is a map of locked seat identifiers and locker user ids.
	Locked map[string]int

	// Reserved is a map with all reserved places and their expiration
	// time. On the expiration time, it can be determined as unreserved.
	Reserved map[string]*time.Time

	// Paid is a map containing all paid seats. There is no need to save
	// other data as a value so we went with the empty struct data type
	// which uses 0 bytes of memory.
	Paid map[string]struct{}
}

// State is a structure containing data about all seats in the whole
// hall (in all rooms). Reserved and Paid seats are got from the database
// on initialization of the program and then the data are received using
// the websocket api from our reservation backend server.
var State = state{
	Locked:   make(map[string]int),
	Reserved: make(map[string]*time.Time),
	Paid:     make(map[string]struct{}),
}

// Reservation is a simplified database model for getting only data
// necessary for initialization and filling up the current states.
type Reservation struct {
	Status  string
	Payment struct {
		Due *time.Time `bson:"due"`
	} `bson:"payment"`
	Places struct {
		ActiveSeat []string `bson:"activeSeat"`
	} `bson:"places"`
}

// SetCurrentStatus gets the current seat statuses from database and
// fills the State variable with them (Reserved and Paid only, Locked
// seats are not stored in database so if the websocket is restarted,
// the data are lost).
func SetCurrentStatus() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// get all reservations for the currently working year
	cursor, err := db.Collection("buy").Find(ctx, bson.M{
		"year": config.Data.Purples.Year,
	})
	if err != nil {
		log.Println("Cannot get current seat status from database :(")
		log.Fatal(err.Error())
	}
	defer cursor.Close(ctx)

	// fill the data to the state object(s)
	for cursor.Next(ctx) {
		var res Reservation
		cursor.Decode(&res)

		// process only waiting and paid reservations, the canceled ones
		// do not have to be processed because their seats are irrelevant
		if res.Status == "waiting" {
			for _, p := range res.Places.ActiveSeat {
				State.Reserved[p] = res.Payment.Due
			}
		} else if res.Status == "paid" {
			for _, p := range res.Places.ActiveSeat {
				State.Paid[p] = struct{}{}
			}
		}
	}
}
