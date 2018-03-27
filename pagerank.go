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
func (in *PageRankIngester) Do() (chan Rank, error) {
	in.Timestamp = time.Now().UTC()
	ch := make(chan Rank)

	go chanFromURLZip(pageRankTop10M, in.Description, ch)

	return ch, nil
}
