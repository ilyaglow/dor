package dor

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	majesticTop1M = "http://downloads.majestic.com/majestic_million.csv"
)

// MajesticCollection is a List implementation which downloads data
// and translates it to LookupMap
//
// More info: https://blog.majestic.com/development/alexa-top-1-million-sites-retired-heres-majestic-million/
type MajesticCollection struct {
	Collection
	resp *http.Response
}

// fetch send request to server with the data
func (f *MajesticCollection) fetch(url string) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	f.resp = r
	return nil
}

// process represents filling the map with response body data
func (f *MajesticCollection) process() error {
	defer f.resp.Body.Close()

	m := make(LookupMap)
	scanner := bufio.NewScanner(f.resp.Body)

	for scanner.Scan() {
		scanner.Text() // skip header
		parts := strings.Split(scanner.Text(), ",")
		pos, d := parts[0], parts[2]
		pint, _ := strconv.ParseInt(pos, 10, 32)
		puint := uint(pint)
		m[d] = puint
	}

	f.Lock()
	f.Map = m
	f.Timestamp = time.Now().UTC()
	f.Unlock()

	return nil
}

// Do implements filling a map with the data from Majestic CSV file
func (f *MajesticCollection) Do() error {
	f.Description = "majestic"
	if err := f.fetch(majesticTop1M); err != nil {
		return err
	}

	if err := f.process(); err != nil {
		return err
	}
	return nil
}
