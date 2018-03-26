package dor

import (
	"bufio"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	majesticTop1M = "http://downloads.majestic.com/majestic_million.csv"
)

// MajesticIngester is a List implementation which downloads data
// and translates it to LookupMap
//
// More info: https://blog.majestic.com/development/alexa-top-1-million-sites-retired-heres-majestic-million/
type MajesticIngester struct {
	IngesterConf
	resp *http.Response
}

// MajesticRank implements Rank interface
type MajesticRank struct {
	GlobalRank     uint   `json:"rank" bson:"rank"`
	TLDRank        uint   `json:"tld_rank" bson:"tld_rank"`
	Domain         string `json:"domain" bson:"domain"`
	TLD            string `json:"tld" bson:"tld"`
	RefSubNets     uint   `json:"ref_sub_nets" bson:"ref_sub_nets"`
	RefIPs         uint   `json:"ref_ips" bson:"ref_ips"`
	IDNDomain      string `json:"idn_domain" bson:"idn_domain"`
	IDNTLD         string `json:"idn_tld" bson:"idn_tld"`
	PrevGlobalRank uint   `json:"prev_global_rank" bson:"prev_global_rank"`
	PrevTLDRank    uint   `json:"prev_tld_rank" bson:"prev_tld_rank"`
	PrevRefSubNets uint   `json:"prev_ref_sub_nets" bson:"prev_ref_sub_nets"`
	PrevRefIPs     uint   `json:"prev_ref_ips" bson:"prev_ref_ips"`
}

// NewMajestic bootstraps MajesticIngester
func NewMajestic() *MajesticIngester {
	return &MajesticIngester{
		IngesterConf: IngesterConf{
			Description: "majestic",
		},
	}
}

// GetDomain is a simple getter for the MajesticRank's domain
func (m *MajesticRank) GetDomain() string { return m.Domain }

// GetRank is a simple getter for the MajesticRank's rank
func (m *MajesticRank) GetRank() uint { return m.GlobalRank }

// fetch send request to server with the data
func (in *MajesticIngester) fetch(url string) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	in.resp = r
	return nil
}

// process represents filling the map with response body data
func (in *MajesticIngester) process(rc chan Rank) {
	defer in.resp.Body.Close()

	scanner := bufio.NewScanner(in.resp.Body)

	for scanner.Scan() {
		scanner.Text() // skip header
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 12 {
			log.Println("majestic: wrong line in a CSV")
			continue
		}

		rc <- &MajesticRank{
			GlobalRank:     strToUint(parts[0]),
			TLDRank:        strToUint(parts[1]),
			Domain:         parts[2],
			TLD:            parts[3],
			RefSubNets:     strToUint(parts[4]),
			RefIPs:         strToUint(parts[5]),
			IDNDomain:      parts[6],
			IDNTLD:         parts[7],
			PrevGlobalRank: strToUint(parts[8]),
			PrevTLDRank:    strToUint(parts[9]),
			PrevRefSubNets: strToUint(parts[10]),
			PrevRefIPs:     strToUint(parts[11]),
		}
	}

	close(rc)
}

// Do implements Ingester interface with the data from Majestic CSV file
func (in *MajesticIngester) Do() (chan Rank, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan Rank)

	if err := in.fetch(majesticTop1M); err != nil {
		return nil, err
	}

	go in.process(ch)

	return ch, nil
}
