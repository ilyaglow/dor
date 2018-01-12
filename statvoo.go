package dor

import (
	"time"
)

const (
	statvooTop1M = "https://statvoo.com/dl/top-1million-sites.csv.zip"
)

// StatvooCollection represents top 1 million websites by statvoo
//
// More info: https://statvoo.com/top/sites
type StatvooCollection struct {
	Collection
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
