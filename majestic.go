/*
Majestic is a List implementation which downloads data
and translates it to LookupMap
*/
package dor

import (
	"bufio"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	majesticTop1M = "http://downloads.majestic.com/majestic_million.csv"
)

type MajesticCollection struct {
	sync.Mutex
	Description string
	Map         LookupMap
	Timestamp   time.Time
	resp        *http.Response
}

func (f *MajesticCollection) fetch(url string) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	f.resp = r
	return nil
}

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

func (f *MajesticCollection) GetTime() time.Time {
	return f.Timestamp
}

func (f *MajesticCollection) GetDesc() string {
	return f.Description
}

func (f *MajesticCollection) Get(d string) (rank uint, presence bool) {
	rank, prs := f.Map[d]
	return rank, prs
}
