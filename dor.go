package dor

import (
	"log"
	"sync"
	"time"
)

// collections is a slice of implemented List structs
var collections = []List{
	&AlexaCollection{},
	&MajesticCollection{},
	&UmbrellaCollection{},
	&StatvooCollection{},
}

// RankResponse represents domain rank response
// that server sends back to client
type RankResponse struct {
	Domain      string    `json:"domain"`
	Rank        uint      `json:"rank"`
	LastUpdate  time.Time `json:"last_update"`
	Description string    `json:"description"`
}

// Response is a find response
type FindResponse struct {
	RequestData string          `json:"data"`
	Hits        []*RankResponse `json:"ranks"`
	Timestamp   time.Time       `json:"timestamp"`
}

// DomainRank represents the top level structure
type DomainRank struct {
	sync.Mutex
	data []List
}

// Fill fills available List interfaces
func (d *DomainRank) Fill() error {
	for _, c := range collections {
		go func(c List) {
			if err := c.Do(); err != nil {
				log.Printf("Failed to enrich %s: %s", c.GetDesc(), err.Error())
			} else {
				log.Printf("%s is done", c.GetDesc())
			}
			d.Lock()
			d.data = append(d.data, c)
			d.Unlock()
		}(c)
	}
	return nil
}

// Update represents a LookupMap updating inside List implementations
func (d *DomainRank) Update() {
	for _, c := range d.data {
		go func() {
			if err := c.Do(); err != nil {
				log.Printf("Failed to update %s: %s", c.GetDesc(), err.Error())
			} else {
				log.Printf("Successfully updated %s", c.GetDesc())
			}
		}()
	}
}

// FillAndUpdate combines filling and updating on a specific duration
func (d *DomainRank) FillAndUpdate(duration time.Duration) error {
	if err := d.Fill(); err != nil {
		return err
	}

	ticker := time.NewTicker(duration)
	go func() {
		<-ticker.C
		d.Update()
	}()

	return nil
}

// Find represents find operation on all Lists available
func (d *DomainRank) Find(domain string) *FindResponse {
	ranks := []*RankResponse{}

	for _, c := range d.data {
		rank, pr := c.Get(domain)
		if pr == false {
			continue
		}
		r := &RankResponse{
			Domain:      domain,
			Rank:        rank,
			LastUpdate:  c.GetTime(),
			Description: c.GetDesc(),
		}
		ranks = append(ranks, r)
	}

	return &FindResponse{
		RequestData: domain,
		Hits:        ranks,
		Timestamp:   time.Now().UTC(),
	}
}
