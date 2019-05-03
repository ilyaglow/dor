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
	log.Printf("%s downloaded successfully", url)
	in.resp = r
	return nil
}

// process represents filling the map with response body data
func (in *MajesticIngester) process(rc chan *Entry) {
	defer in.resp.Body.Close()

	scanner := bufio.NewScanner(in.resp.Body)
	scanner.Text() // skip header

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 12 {
			log.Println("majestic: wrong line in a CSV")
			continue
		}

		rc <- &Entry{
			Rank:    strToUint(parts[0]),
			Domain:  parts[2],
			RawData: []byte(line),
		}
	}

	close(rc)
}

// Do implements Ingester interface with the data from Majestic CSV file
func (in *MajesticIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	if err := in.fetch(majesticTop1M); err != nil {
		return nil, err
	}

	go in.process(ch)

	return ch, nil
}
