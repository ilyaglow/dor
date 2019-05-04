package dor

import (
	"sync"
	"time"
)

// LookupMap represents map with domain - rank pairs
type LookupMap map[string]uint

// MemoryCollection is a struct that is capable to hold data
type memoryCollection struct {
	sync.Mutex
	Map         LookupMap
	Description string
	Timestamp   time.Time
}

// MemoryStorage implements Storage interface as in-memory storage
type MemoryStorage struct {
	Maps map[string]*memoryCollection
}

func (mc *memoryCollection) get(d string) (rank uint, presence bool) {
	rank, pr := mc.Map[d]
	return rank, pr
}

// Get implements Get method of the Storage interface
func (ms *MemoryStorage) Get(d string, sources ...string) ([]*Entry, error) {
	var ranks []*Entry
	for k := range ms.Maps {
		if len(sources) > 0 && !sliceContains(sources, ms.Maps[k].Description) {
			continue
		}

		rank, prs := ms.Maps[k].get(d)
		if prs != true {
			continue
		}

		r := &Entry{
			Domain: d,
			Rank:   rank,
			Date:   ms.Maps[k].Timestamp,
			Source: ms.Maps[k].Description,
		}

		ranks = append(ranks, r)
	}

	return ranks, nil
}

// GetMore is not supported for the memory storage
func (ms *MemoryStorage) GetMore(d string, lps int, sources ...string) ([]*Entry, error) {
	return ms.Get(d, sources...)
}

// Put implements Put method of the Storage interface
func (ms *MemoryStorage) Put(c <-chan *Entry, s string, t time.Time) error {
	lm := make(LookupMap)

	for r := range c {
		lm[r.Domain] = r.Rank
	}

	if _, ok := ms.Maps[s]; ok {
		ms.Maps[s].Map = lm
		return nil
	}

	ms.Maps[s] = &memoryCollection{
		Map:         lm,
		Description: s,
		Timestamp:   t,
	}

	return nil
}

func sliceContains(ss []string, s string) bool {
	for i := range ss {
		if ss[i] == s {
			return true
		}
	}
	return false
}
