package nhc

import (
	"encoding/json"
)

const (
	RegisterCMD     = "{\"cmd\":\"startevents\"}"
	ListActions     = "{\"cmd\":\"listactions\"}"
	ListLocations   = "{\"cmd\":\"listlocations\"}"
	ListEnergies    = "{\"cmd\":\"listenergy\"}"
	ListThermostats = "{\"cmd\":\"listthermostat\"}"
)

// Message generic struct to hold nhc messages
// used to identify the message type before futher parsing
type Message struct {
	Cmd   string `json:"cmd"`
	Event string `json:"event"`
	//Data []NhcAction `json:"data"`
	//Data []interface{} `json:"data"`
	Data json.RawMessage
}

// Action holds one individual nhc action (equipment)
type Action struct {
	ID       int
	Name     string
	Type     int
	Location int
	Value1   int
	Value2   int
	Value3   int
}

// Event holds an individual event
type Event struct {
	ID    int `json:"id"`
	Value int `json:"value1"`
}

// Location holds one nhc location
type Location struct {
	ID   int
	Name string
}

// SimpleCmd type holding a nhc command
type SimpleCmd struct {
	Cmd   string `json:"cmd"`
	ID    int    `json:"id"`
	Value int    `json:"value1"`
}

// Stringify return the string version of SimpleCmd
func (sc SimpleCmd) Stringify() string {
	tmp, _ := json.Marshal(sc)
	return string(tmp)
}
