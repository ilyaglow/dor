package dor

import "time"

const (
	quantcastTop1M = "https://ak.quantcast.com/quantcast-top-sites.zip"
)

// QuantcastIngester represents Ingester implementation for Quantcast Top 1
// Million websites.
type QuantcastIngester struct {
	IngesterConf
}

// NewQuantcast bootstraps QuantcastIngester.
func NewQuantcast() *QuantcastIngester {
	return &QuantcastIngester{
		IngesterConf: IngesterConf{
			Description: "quantcast",
		},
	}

}

// Do gets the data from Quantcast Top 1M txt file.
func (in *QuantcastIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	go chanFromURLZip(quantcastTop1M, in.Description, ch, "	", 7)

	return ch, nil
}
