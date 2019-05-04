package dor

import (
	"time"
)

const (
	alexaTop1M = "http://s3.amazonaws.com/alexa-static/top-1m.csv.zip"
)

// AlexaIngester represents Ingester implementation for Alexa Top 1 Million websites
type AlexaIngester struct {
	IngesterConf
}

// NewAlexa bootstraps AlexaIngester
func NewAlexa() *AlexaIngester {
	return &AlexaIngester{
		IngesterConf: IngesterConf{
			Description: "alexa",
		},
	}

}

// Do implements Ingester Do func with the data from Alexa Top 1M CSV file
func (in *AlexaIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	go chanFromURLZip(alexaTop1M, in.Description, ch, ",", 0)

	return ch, nil
}
