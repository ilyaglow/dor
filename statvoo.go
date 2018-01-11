package dor

import (
	"sync"
	"time"
)

const (
	statvooTop1M = "https://statvoo.com/dl/top-1million-sites.csv.zip"
)

type StatvooCollection struct {
	sync.Mutex
	Description string
	Map         LookupMap
	Timestamp   time.Time
}

func (f *StatvooCollection) Do() error {
	f.Description = "statvoo"

	m, err := mapFromURLZip(statvooTop1M, f.Description)
	if err != nil {
		return err
	}
	f.Lock()
	f.Map = m
	f.Timestamp = time.Now().UTC()
	f.Unlock()

	return nil
}

func (f *StatvooCollection) GetTime() time.Time {
	return f.Timestamp
}

func (f *StatvooCollection) GetDesc() string {
	return f.Description
}

func (f *StatvooCollection) Get(d string) (rank uint, presence bool) {
	rank, prs := f.Map[d]
	return rank, prs
}
