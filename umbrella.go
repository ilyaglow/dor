package dor

import (
	"time"
)

const (
	umbrellaTop1M = "http://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip"
)

// UmbrellaIngester represents Ingester implementation for OpenDNS Umbrella Top 1M domains
//
// More info: https://umbrella.cisco.com/blog/2016/12/14/cisco-umbrella-1-million/
type UmbrellaIngester struct {
	IngesterConf
}

// NewUmbrella bootstraps UmbrellaIngester
func NewUmbrella() *UmbrellaIngester {
	return &UmbrellaIngester{
		IngesterConf: IngesterConf{
			Description: "umbrella",
		},
	}
}

// Do implements Ingester Do func with the data from OpenDNS
func (in *UmbrellaIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	go chanFromURLZip(umbrellaTop1M, in.Description, ch, ",", 0)

	return ch, nil
}
