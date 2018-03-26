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
func (ms *MemoryStorage) Get(d string, sources ...string) ([]Rank, error) {
	var ranks []Rank
	for k := range ms.Maps {
		rank, prs := ms.Maps[k].get(d)
		if prs != true {
			continue
		}

		r := ExtendedRank{
			Domain:     d,
			Rank:       rank,
			LastUpdate: ms.Maps[k].Timestamp,
			Source:     ms.Maps[k].Description,
		}

		ranks = append(ranks, r)
	}

	return ranks, nil
}

// Put implements Put method of the Storage interface
func (ms *MemoryStorage) Put(c <-chan Rank, s string, t time.Time) error {
	for r := range c {
		if _, ok := ms.Maps[s]; !ok {
			ms.Maps[s] = &memoryCollection{
				Map:         make(LookupMap),
				Description: s,
				Timestamp:   t,
			}
		}

		ms.Maps[s].Map[r.GetDomain()] = r.GetRank()
	}

	return nil
}
