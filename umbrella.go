package dor

import (
	"time"
)

const (
	umbrellaTop1M = "http://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip"
)

// UmbrellaCollection represents List implementation for OpenDNS Umbrella Top 1M domains
//
// More info: https://umbrella.cisco.com/blog/2016/12/14/cisco-umbrella-1-million/
type UmbrellaCollection struct {
	Collection
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
