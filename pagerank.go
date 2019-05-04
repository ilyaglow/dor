package dor

import "time"

const pageRankTop10M = "https://www.domcop.com/files/top/top10milliondomains.csv.zip"

// PageRankIngester represents Ingester implementation for Domcop PageRank top 10M domains
type PageRankIngester struct {
	IngesterConf
}

// NewPageRank bootstraps PageRankIngester
func NewPageRank() *PageRankIngester {
	return &PageRankIngester{
		IngesterConf: IngesterConf{
			Description: "pagerank",
		},
	}
}

// Do implements Ingester Do func with the data from DomCop
func (in *PageRankIngester) Do() (chan *Entry, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan *Entry)

	go chanFromURLZip(pageRankTop10M, in.Description, ch, ",", 0)

	return ch, nil
}
