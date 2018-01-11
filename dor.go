package dor

import (
	"log"
	"time"
)

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
	data []List
}

// Fill fills available List interfaces
func (d *DomainRank) Fill() error {
	alexa := &AlexaCollection{}
	if err := alexa.Do(); err != nil {
		return err
	}
	log.Println("Alexa is done")
	d.data = append(d.data, alexa)

	majestic := &MajesticCollection{}
	if err := majestic.Do(); err != nil {
		return err
	}
	log.Println("Majestic is done")
	d.data = append(d.data, majestic)

	umbrella := &UmbrellaCollection{}
	if err := umbrella.Do(); err != nil {
		return err
	}
	log.Println("Umbrella is done")
	d.data = append(d.data, umbrella)

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
