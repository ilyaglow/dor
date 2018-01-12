package dor

import (
	"sync"
	"time"
)

// collections is a slice of implemented List structs
var collections = []List{
	&AlexaCollection{},
	&MajesticCollection{},
	&UmbrellaCollection{},
	&StatvooCollection{},
}

// Collection represents a top popular domains provider.
//
// Implemented collections now are:
//	- Alexa Top 1 Million
//	- Majestic Top 1 Million
//	- Umbrella Top 1 Million
//	- Statvoo Top 1 Million
type Collection struct {
	sync.Mutex
	Description string
	Map         LookupMap
	Timestamp   time.Time
}

// Get returns query result from a map.
func (c *Collection) Get(d string) (rank uint, presence bool) {
	rank, prs := c.Map[d]
	return rank, prs
}

// GetDesc is a simple getter for a collection's description
func (c *Collection) GetDesc() string {
	return c.Description
}

// GetTime is a simple getter for a collection's timestamp
func (c *Collection) GetTime() time.Time {
	return c.Timestamp
}
