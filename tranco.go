package dor

import (
	"time"
)

const (
	trancoTop1M = "https://tranco-list.eu/top-1m.csv.zip"
)

// TrancoIngester represents Ingester implementation for Tranco Top 1 Million
// websites.
// About: https://tranco-list.eu/
type TrancoIngester struct {
	IngesterConf
}

// NewTranco bootstraps TrancoIngester
func NewTranco() *TrancoIngester {
	return &TrancoIngester{
		IngesterConf: IngesterConf{
			Description: "tranco",
		},
	}
}

// Do implements Ingester Do func with the data from Tranco Top 1M CSV file
func (in *TrancoIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	go chanFromURLZip(trancoTop1M, in.Description, ch)

	return ch, nil
}
