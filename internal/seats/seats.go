package seats

import "time"

// Request is a structure which represents a payload from websocket
// informing backend about one specific action (locked or unlocked seat)
type Request struct {
	Action              string `json:"action"`
	Seat                string `json:"seat"`
	AuthorizationString string `json:"authorizationString"`
}

// State is a total state for all seats in the whole hall (in all rooms)
type State struct {
	Reserved []string `json:"reserved"`
	Paid     []string `json:"paid"`
	Locked   []string `json:"locked"`
}

// Locked is a map of locked seat identifiers and locker user ids
var Locked = make(map[string]int)

// Reserved is a map with all reserved places and their expiration
// time. On the expiration time, it can be determined as unreserved.
var Reserved = make(map[string]*time.Time)

// Paid is a map containing all paid seats
var Paid = make(map[string]bool)
