package dor

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// // Source is an interface for collections.
// type Source interface {
// 	Do() error               // fills the Storage
// 	Get(string) (uint, bool) // gets rank from the Storage
// 	GetDesc() string         // description getter
// 	GetTime() time.Time      // time getter
// }

// Ingester fetches data and uploads it to the Storage
type Ingester interface {
	Do() (chan Rank, error) // returns a channel for consumers
	GetDesc() string        // simple getter for the source
}

// Storage represents an interface to store and query ranks.
type Storage interface {
	Put(<-chan Rank, string, time.Time) error                          // Put is usually a bulk inserter from the channel that works in a goroutine, second argument is a Source of the data and third is the last update time
	Get(domain string, sources ...string) ([]Rank, error)              // Get is a simple getter for the latest rank of the domain in a particular domain rank provider
	GetMore(domain string, lps int, sources ...string) ([]Rank, error) // GetAll is a getter that retreives historical data on the domain limited by lps (limit per source)
}

// Rank is an interface for different ranking systems
type Rank interface {
	GetDomain() string
	GetRank() uint
}

// SimpleRank is a simple domain rank structure.
type SimpleRank struct {
	Domain string `json:"domain" db:"domain" bson:"domain"`
	Rank   uint   `json:"rank" db:"rank" bson:"rank"`
}

// ExtendedRank is a SimpleRank with extended fields
type ExtendedRank struct {
	Domain     string    `json:"domain" db:"domain" bson:"domain"`
	Rank       uint      `json:"rank" db:"rank" bson:"rank"`
	LastUpdate time.Time `json:"last_update" bson:"last_update"`
	Source     string    `json:"source" bson:"source"`
}

// FindResponse is a find request response.
type FindResponse struct {
	RequestData string    `json:"data"`
	Hits        []Rank    `json:"ranks"`
	Timestamp   time.Time `json:"timestamp"`
}

// App represents Dor configuration options
type App struct {
	Ingesters []Ingester
	Storage   Storage
	Keep      bool
}

// GetDomain is a simple getter for a Domain
func (s SimpleRank) GetDomain() string { return s.Domain }

// GetRank is a simple getter for a Rank
func (s SimpleRank) GetRank() uint { return s.Rank }

// GetDomain is a simple getter for a Domain
func (s ExtendedRank) GetDomain() string { return s.Domain }

// GetRank is a simple getter for a Rank
func (s ExtendedRank) GetRank() uint { return s.Rank }

// New bootstraps App struct.
//	stn - storage name
//	stl - storage location string
//	keep - keep new data or overwrite old one (always false for MemoryStorage)
func New(stn string, stl string, keep bool) (*App, error) {
	var s Storage
	var err error
	switch stn {
	case "memory":
		s = &MemoryStorage{make(map[string]*memoryCollection)}
	case "mongodb":
		s, err = NewMongoStorage(stl, "dor", "ranks", 50000, 5, keep)
		if err != nil {
			return nil, err
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
				log.Printf("Failed to enrich %s: %s", ing.GetDesc(), err.Error())
				return
			}

			if err := d.Storage.Put(ch, ing.GetDesc(), time.Now().UTC()); err != nil {
				log.Printf("Failed to insert data to the storage %s: %s", ing.GetDesc(), err.Error())
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
	go func() {
		<-ticker.C
		d.Fill()
	}()

	return nil
}

// Find represents find operation on the storage available
func (d *App) Find(domain string, sources ...string) (*FindResponse, error) {
	var ranks []Rank
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
