package dor

import (
	"time"
)

const (
	alexaTop1M = "http://s3.amazonaws.com/alexa-static/top-1m.csv.zip"
)

// AlexaCollection represents List implementation for Alexa Top 1 Million websites
type AlexaCollection struct {
	Collection
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
