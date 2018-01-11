package dor

import (
	"sync"
	"time"
)

const (
	statvooTop1M = "https://statvoo.com/dl/top-1million-sites.csv.zip"
)

// StatvooCollection represents top 1 million websites by statvoo
//
// More info: https://statvoo.com/top/sites
type StatvooCollection struct {
	sync.Mutex
	Description string
	Map         LookupMap
	Timestamp   time.Time
}

// Do implements filling a LookupMap from the data source
func (f *StatvooCollection) Do() error {
	f.Description = "statvoo"

	m, err := mapFromURLZip(statvooTop1M, f.Description)
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
func (f *StatvooCollection) GetTime() time.Time {
	return f.Timestamp
}

// GetDesc represents a collection description
func (f *StatvooCollection) GetDesc() string {
	return f.Description
}

// Get represents a specific domain rank finding
func (f *StatvooCollection) Get(d string) (rank uint, presence bool) {
	rank, prs := f.Map[d]
	return rank, prs
}
