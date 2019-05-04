package dor

import (
	"time"
)

const (
	statvooTop1M = "https://statvoo.com/dl/top-1million-sites.csv.zip"
)

// StatvooIngester represents top 1 million websites by statvoo
//
// More info: https://statvoo.com/top/sites
type StatvooIngester struct {
	IngesterConf
}

// NewStatvoo boostraps StatvooIngester
func NewStatvoo() *StatvooIngester {
	return &StatvooIngester{
		IngesterConf: IngesterConf{
			Description: "statvoo",
		},
	}
}

// Do implements Ingester Do func with the data from Statvoo Top 1M
func (in *StatvooIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	go chanFromURLZip(statvooTop1M, in.Description, ch, ",", 0)

	return ch, nil
}
