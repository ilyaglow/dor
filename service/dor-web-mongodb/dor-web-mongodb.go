package main

import (
	"flag"
	"fmt"

	"github.com/ilyaglow/dor"
	"github.com/ilyaglow/dor/web"
)

func main() {
	bindAddr := flag.String("host", "127.0.0.1", "IP-address to bind")
	bindPort := flag.String("port", "8080", "Port to bind")
	flag.Parse()
	hp := fmt.Sprintf("%s:%s", *bindAddr, *bindPort)

	d, err := dor.New("mongodb", "", false)
	if err != nil {
		panic(err)
	}

	dorweb.Serve(hp, d)
}
