package dor

import (
	"sync"
	"time"
)

const (
	umbrellaTop1M = "http://s3-us-west-1.amazonaws.com/umbrella-static/top-1m.csv.zip"
)

type UmbrellaCollection struct {
	sync.Mutex
	Description string
	Map         LookupMap
	Timestamp   time.Time
}

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

func (f *UmbrellaCollection) GetTime() time.Time {
	return f.Timestamp
}

func (f *UmbrellaCollection) GetDesc() string {
	return f.Description
}

func (f *UmbrellaCollection) Get(d string) (rank uint, presence bool) {
	rank, prs := f.Map[d]
	return rank, prs
}
