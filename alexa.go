package dor

import (
	"archive/zip"
	"io"
	"sync"
	"time"
)

const (
	alexaTop1M = "http://s3.amazonaws.com/alexa-static/top-1m.csv.zip"
)

type AlexaCollection struct {
	sync.Mutex
	Description string
	tmpFile     *zip.ReadCloser
	tmpFileName string
	csvFile     *io.ReadCloser
	Map         LookupMap
	Timestamp   time.Time
}

func (f *AlexaCollection) Do() error {
	f.Description = "alexa"

	m, err := mapFromURLZip(alexaTop1M, f.Description)
	if err != nil {
		return err
	}
	f.Lock()
	f.Map = m
	f.Timestamp = time.Now().UTC()
	f.Unlock()

	return nil
}

func (f *AlexaCollection) GetTime() time.Time {
	return f.Timestamp
}

func (f *AlexaCollection) GetDesc() string {
	return f.Description
}

func (f *AlexaCollection) Get(d string) (rank uint, presence bool) {
	rank, prs := f.Map[d]
	return rank, prs
}
