package main

import (
	"flag"
	"log"
	"time"

	"github.com/ilyaglow/dor"
	"github.com/ilyaglow/dor/web"
)

const duration = time.Hour * 24

func main() {
	bindAddr := flag.String("listen", "127.0.0.1:8080", "Listen address and port to bind")
	flag.Parse()
	d, err := dor.New("memory", "", false)
	if err != nil {
		panic(err)
	}

	go func() {
		if err := d.FillByTimer(duration); err != nil {
			log.Fatal(err)
		}
	}()

	dorweb.Serve(*bindAddr, d)
}
