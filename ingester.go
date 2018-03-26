package dor

import (
	"sync"
	"time"
)

// ingesters is a slice of implemented Ingester structs
var ingesters = []Ingester{
	NewAlexa(),
	NewStatvoo(),
	NewUmbrella(),
	NewMajestic(),
}

// IngesterConf represents a top popular domains provider configuration.
//
// Implemented ingesters by now are:
//	- Alexa Top 1 Million
//	- Majestic Top 1 Million
//	- Umbrella Top 1 Million
//	- Statvoo Top 1 Million
type IngesterConf struct {
	sync.Mutex
	Description string
	Timestamp   time.Time
}

// GetDesc is a simple getter for a collection's description
func (in *IngesterConf) GetDesc() string {
	return in.Description
}
