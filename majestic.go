package dor

import (
	"bufio"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	majesticTop1M = "http://downloads.majestic.com/majestic_million.csv"
)

// MajesticIngester is a List implementation which downloads data
// and translates it to LookupMap
//
// More info: https://blog.majestic.com/development/alexa-top-1-million-sites-retired-heres-majestic-million/
type MajesticIngester struct {
	IngesterConf
	resp *http.Response
}

// NewMajestic bootstraps MajesticIngester
func NewMajestic() *MajesticIngester {
	return &MajesticIngester{
		IngesterConf: IngesterConf{
			Description: "majestic",
		},
	}
}

// fetch send request to server with the data
func (in *MajesticIngester) fetch(url string) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	log.Printf("%s downloaded successfully", url)
	in.resp = r
	return nil
}

// process represents filling the map with response body data
func (in *MajesticIngester) process(rc chan *Entry) {
	defer in.resp.Body.Close()

	scanner := bufio.NewScanner(in.resp.Body)

	var i int
	for scanner.Scan() {
		line := scanner.Text()
		if i < 1 {
			i++
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 12 {
			log.Println("majestic: wrong line in a CSV")
			continue
		}

		rc <- &Entry{
			Rank:    strToUint(parts[0]),
			Domain:  parts[2],
			RawData: line,
		}
	}

	close(rc)
}

// Do implements Ingester interface with the data from Majestic CSV file
func (in *MajesticIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	if err := in.fetch(majesticTop1M); err != nil {
		return nil, err
	}

	go in.process(ch)

	return ch, nil
}
