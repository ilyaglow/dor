package dor

import (
	"sync"
	"time"
)

const (
	umbrellaTop1M = "http://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip"
)

// UmbrellaCollection represents List implementation for OpenDNS Umbrella Top 1M domains
//
// More info: https://umbrella.cisco.com/blog/2016/12/14/cisco-umbrella-1-million/
type UmbrellaCollection struct {
	sync.Mutex
	Description string
	Map         LookupMap
	Timestamp   time.Time
}

// Do implements filling a map with the data from OpenDNS
func (f *UmbrellaCollection) Do() error {
	f.Description = "umbrella"

	m, err := mapFromURLZip(umbrellaTop1M, f.Description)
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
func (f *UmbrellaCollection) GetTime() time.Time {
	return f.Timestamp
}

// GetDesc represents a collection description
func (f *UmbrellaCollection) GetDesc() string {
	return f.Description
}

// Get represents a specific domain rank finding
func (f *UmbrellaCollection) Get(d string) (rank uint, presence bool) {
	rank, prs := f.Map[d]
	return rank, prs
}
