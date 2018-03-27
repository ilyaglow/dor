package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/ilyaglow/dor"
	"github.com/ilyaglow/dor/web"
)

const duration = time.Hour * 24

func main() {
	bindAddr := flag.String("host", "127.0.0.1", "IP-address to bind")
	bindPort := flag.String("port", "8080", "Port to bind")
	flag.Parse()
	hp := fmt.Sprintf("%s:%s", *bindAddr, *bindPort)

	d, err := dor.New("mongodb", "", false)
	if err != nil {
		panic(err)
	}

	if err := d.FillByTimer(duration); err != nil {
		log.Fatal(err)
	}

	dorweb.Serve(hp, d)
}
