package seats

// Request is a structure which represents a payload from websocket
// informing backend about one specific action (locked or unlocked seat)
type Request struct {
	Action string `json:"action"`
	Seat   string `json:"seat"`
}

// State is a total state for all seats in the whole hall (in all rooms)
type State struct {
	Reserved []string `json:"reserved"`
	Paid     []string `json:"paid"`
	Locked   []string `json:"locked"`
}

// Locked is a array of locked seat identifiers
var Locked = make(map[string]int)
