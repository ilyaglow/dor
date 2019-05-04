package main

import (
	"flag"
	"log"
	"os"

	"github.com/ilyaglow/dor"
	"github.com/ilyaglow/dor/web"
	"github.com/peterbourgon/ff"
)

func main() {
	fs := flag.NewFlagSet("DOR", flag.ExitOnError)
	var (
		storage  = fs.String("storage", "clickhouse", "storage type")
		location = fs.String("storage-url", "", "url of the storage")
		listen   = fs.String("listen-addr", ":8080", "listen address")
	)
	ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("DOR"))

	d, err := dor.New(*storage, *location, false)
	if err != nil {
		panic(err)
	}

	log.Fatal(dorweb.Serve(*listen, d))
}
