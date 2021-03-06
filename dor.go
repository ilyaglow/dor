package dor

import (
	"fmt"
	"log"
	"sync"
	"time"
)

const (
	batchSize = 50000
	tblName   = "ranks"
)

// DefaultTTL for records in days.
var DefaultTTL = 30

// Storage represents an interface to store and query ranks.
type Storage interface {
	Put(<-chan *Entry, string, time.Time) error             // Put is usually a bulk inserter from the channel that works in a goroutine, second argument is a Source of the data and third is the last update time.
	Get(domain string, sources ...string) ([]*Entry, error) // Get is a simple getter for the latest rank of the domain in a particular domain rank provider or all of them if nothing selected.
}

// Entry is a SimpleRank with extended fields
type Entry struct {
	Domain  string    `json:"domain" db:"domain" bson:"domain"`
	Rank    uint32    `json:"rank" db:"rank" bson:"rank"`
	Date    time.Time `json:"date" bson:"date"`
	Source  string    `json:"source" bson:"source"`
	RawData string    `json:"raw" bson:"raw"`
}

// FindResponse is a find request response.
type FindResponse struct {
	RequestData string    `json:"data"`
	Hits        []*Entry  `json:"ranks"`
	Timestamp   time.Time `json:"timestamp"`
}

// App represents Dor configuration options
type App struct {
	Ingesters []Ingester
	Storage   Storage
	Keep      bool
}

// New bootstraps App struct.
//	stn - storage name
//	stl - storage location string
//	keep - keep new data or overwrite old one (always false for MemoryStorage)
func New(stn string, stl string, keep bool) (*App, error) {
	var (
		s   Storage
		err error
	)
	switch stn {
	case "clickhouse":
		s, err = NewClickhouseStorage(stl, tblName, batchSize)
		if err != nil {
			return nil, fmt.Errorf("new clickhouse storage: %w", err)
		}
	case "memory":
		s = &MemoryStorage{make(map[string]*memoryCollection)}
	case "mongodb":
		s, err = NewMongoStorage(stl, "dor", tblName, batchSize, 5, keep)
		if err != nil {
			return nil, fmt.Errorf("new mongo storage: %w", err)
		}
	default:
		return nil, fmt.Errorf("%s storage is not implemented", stn)
	}

	return &App{
		Ingesters: ingesters,
		Storage:   s,
		Keep:      keep,
	}, nil
}

// Fill fills available Ingester interfaces.
func (d *App) Fill() error {
	var wg sync.WaitGroup
	wg.Add(len(d.Ingesters))

	for _, ing := range d.Ingesters {
		go func(ing Ingester) {
			defer wg.Done()

			ch, err := ing.Do()
			if err != nil {
				log.Printf("failed to enrich %s: %s", ing.GetDesc(), err.Error())
				return
			}

			if err := d.Storage.Put(ch, ing.GetDesc(), time.Now().UTC()); err != nil {
				log.Printf("failed to insert data to the storage %s: %s", ing.GetDesc(), err.Error())
				return
			}

			log.Printf("%s is done", ing.GetDesc())
		}(ing)
	}
	wg.Wait()
	return nil
}

// FillByTimer combines filling and updating on a specific duration
func (d *App) FillByTimer(duration time.Duration) error {
	if err := d.Fill(); err != nil {
		return err
	}

	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	for {
		<-ticker.C
		if err := d.Fill(); err != nil {
			return err
		}
	}
}

// Find represents find operation on the storage available
func (d *App) Find(domain string, sources ...string) (*FindResponse, error) {
	var ranks []*Entry
	var ings []string
	for i := range d.Ingesters {
		ings = append(ings, d.Ingesters[i].GetDesc())
	}

	if len(sources) == 0 {
		sources = ings
	}

	ranks, err := d.Storage.Get(domain, sources...)
	if err != nil {
		return nil, err
	}

	return &FindResponse{
		RequestData: domain,
		Hits:        ranks,
		Timestamp:   time.Now().UTC(),
	}, nil
}
