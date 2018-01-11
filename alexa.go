package dor

import (
	"sync"
	"time"
)

const (
	alexaTop1M = "http://s3.amazonaws.com/alexa-static/top-1m.csv.zip"
)

// AlexaCollection represents List implementation for Alexa Top 1 Million websites
type AlexaCollection struct {
	sync.Mutex
	Description string
	Map         LookupMap
	Timestamp   time.Time
}

// Do implements filling a map with the data from Alexa Top 1M CSV file
func (f *AlexaCollection) Do() error {
	f.Description = "alexa"

	m, err := mapFromURLZip(alexaTop1M, f.Description)
	if err != nil {
		return err
	}
	f.Lock()
	f.Map = m
	f.Timestamp = time.Now().UTC()
	f.Unlock()

	return nil
}

// GetTime represents a collection last updated timestamp
func (f *AlexaCollection) GetTime() time.Time {
	return f.Timestamp
}

// GetDesc represents a collection description
func (f *AlexaCollection) GetDesc() string {
	return f.Description
}

// Get represents a specific domain rank finding
func (f *AlexaCollection) Get(d string) (rank uint, presence bool) {
	rank, prs := f.Map[d]
	return rank, prs
}
